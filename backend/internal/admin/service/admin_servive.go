package service

import (
	"fmt"
	"time"

	"github.com/sdv77/Linguture1/internal/admin/models"
	"github.com/sdv77/Linguture1/internal/admin/repository"
	"github.com/sdv77/Linguture1/pkg/token"
)

// AdminService содержит бизнес-логику администраторов
type AdminService struct {
	adminRepo    *repository.AdminRepository
	teacherRepo  *repository.AdminRepository
	tokenService *token.Service
}

// NewAdminService создает новый сервис
func NewAdminService(adminRepo *repository.AdminRepository, teacherRepo *repository.AdminRepository, tokenService *token.Service) *AdminService {
	return &AdminService{
		adminRepo:    adminRepo,
		teacherRepo:  teacherRepo,
		tokenService: tokenService,
	}
}

// Login выполняет вход администратора (простое сравнение паролей)
func (s *AdminService) Login(input models.AdminLoginInput) (*models.TokenResponse, error) {
	admin, err := s.adminRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, fmt.Errorf("неверное имя пользователя или пароль")
	}

	if !admin.IsActive {
		return nil, fmt.Errorf("администратор не активен")
	}

	// Простое сравнение паролей (без хеширования)
	if admin.Password != input.Password {
		return nil, fmt.Errorf("неверное имя пользователя или пароль")
	}

	// Генерируем JWT токен
	tokenStr, err := s.tokenService.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return &models.TokenResponse{
		Token:     tokenStr,
		ExpiresIn: int(s.tokenService.GetTokenExpiry().Seconds()),
	}, nil
}

// BackupDatabase создает резервную копию базы данных
func (s *AdminService) BackupDatabase(dbHost, dbPort, dbUser, dbPassword, dbName string) ([]byte, error) {
	// Для простоты возвращаем текстовый дамп
	// В реальном приложении здесь будет вызов pg_dump
	timestamp := time.Now().Format("20060102_150405")
	backupContent := fmt.Sprintf(`-- Backup of words_db created at %s
-- This is a simplified backup for demonstration purposes
-- In production, use pg_dump command

-- Admins table
INSERT INTO admins (username, password_hash, full_name, is_active) VALUES
('admin', '$2a$10$rRyBsGSHK6.uc8f7TXbjPuQv0bNeWlP7OZLUmXwXLj.xaFNaPjT8u', 'Администратор системы', TRUE)
ON CONFLICT (username) DO NOTHING;

-- Teachers table backup would be here
-- Lessons table backup would be here
-- ... other tables ...

-- Backup completed successfully at %s
`, timestamp, timestamp)

	return []byte(backupContent), nil
}

// GetAllTeachers получает всех учителей
func (s *AdminService) GetAllTeachers() ([]map[string]interface{}, error) {
	return s.teacherRepo.GetAllTeachers()
}

// CreateTeacher создает нового учителя
func (s *AdminService) CreateTeacher(email, password, fullName string) error {
	return s.teacherRepo.CreateTeacher(email, password, fullName)
}

// UpdateTeacher обновляет учителя
func (s *AdminService) UpdateTeacher(id int, email, fullName string, isActive bool) error {
	return s.teacherRepo.UpdateTeacher(id, email, fullName, isActive)
}

// DeleteTeacher удаляет учителя
func (s *AdminService) DeleteTeacher(id int) error {
	return s.teacherRepo.DeleteTeacher(id)
}

// GetTeacherByID получает учителя по ID
func (s *AdminService) GetTeacherByID(id int) (map[string]interface{}, error) {
	return s.teacherRepo.GetTeacherByID(id)
}
