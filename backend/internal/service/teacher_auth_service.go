package service

import (
	"fmt"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
	"github.com/sdv77/Linguture1/pkg/token"
)

// TeacherAuthService содержит логику аутентификации учителей
type TeacherAuthService struct {
	teacherRepo *repository.TeacherRepository
	tokenSvc    *token.Service
}

// NewTeacherAuthService создает новый сервис
func NewTeacherAuthService(teacherRepo *repository.TeacherRepository, tokenSvc *token.Service) *TeacherAuthService {
	return &TeacherAuthService{
		teacherRepo: teacherRepo,
		tokenSvc:    tokenSvc,
	}
}

// Login выполняет вход учителя (простая проверка пароля)
func (s *TeacherAuthService) Login(input models.TeacherLoginInput) (*models.TokenResponse, error) {
	// Находим учителя по email
	teacher, err := s.teacherRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Проверяем активен ли учитель
	if !teacher.IsActive {
		return nil, fmt.Errorf("учитель не активен")
	}

	// Проверяем пароль (простое сравнение)
	if teacher.Password != input.Password {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Генерируем JWT токен
	tokenStr, err := s.tokenSvc.GenerateToken(teacher.ID, teacher.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return &models.TokenResponse{
		Token:     tokenStr,
		ExpiresIn: int(s.tokenSvc.GetTokenExpiry().Seconds()),
	}, nil
}

// GetTeacherByID получает учителя по ID
func (s *TeacherAuthService) GetTeacherByID(teacherID int) (*models.TeacherResponse, error) {
	teacher, err := s.teacherRepo.GetByID(teacherID)
	if err != nil {
		return nil, err
	}

	return &models.TeacherResponse{
		ID:       teacher.ID,
		Email:    teacher.Email,
		FullName: teacher.FullName,
	}, nil
}
