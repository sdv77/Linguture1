package repository

import (
	"database/sql"
	"fmt"

	"github.com/sdv77/Linguture1/internal/models"
)

// TeacherRepository работает с таблицей учителей
type TeacherRepository struct {
	db *sql.DB
}

// NewTeacherRepository создает новый репозиторий
func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// FindByEmail находит учителя по email
func (r *TeacherRepository) FindByEmail(email string) (*models.Teacher, error) {
	query := `SELECT id, email, password, full_name, is_active, created_at, updated_at 
              FROM teachers WHERE email = $1`

	var teacher models.Teacher
	err := r.db.QueryRow(query, email).Scan(
		&teacher.ID,
		&teacher.Email,
		&teacher.Password, // Теперь получаем пароль напрямую
		&teacher.FullName,
		&teacher.IsActive,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("учитель не найден")
		}
		return nil, fmt.Errorf("ошибка получения учителя: %w", err)
	}

	return &teacher, nil
}

// GetByID получает учителя по ID
func (r *TeacherRepository) GetByID(id int) (*models.Teacher, error) {
	query := `SELECT id, email, full_name, is_active, created_at, updated_at 
              FROM teachers WHERE id = $1`

	var teacher models.Teacher
	err := r.db.QueryRow(query, id).Scan(
		&teacher.ID,
		&teacher.Email,
		&teacher.FullName,
		&teacher.IsActive,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("учитель не найден")
		}
		return nil, fmt.Errorf("ошибка получения учителя: %w", err)
	}

	return &teacher, nil
}
