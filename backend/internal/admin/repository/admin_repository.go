package repository

import (
	"database/sql"
	"fmt"

	"github.com/sdv77/Linguture1/internal/admin/models"
)

// AdminRepository работает с таблицей администраторов
type AdminRepository struct {
	db *sql.DB
}

// NewAdminRepository создает новый репозиторий
func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByUsername находит администратора по имени пользователя
func (r *AdminRepository) FindByUsername(username string) (*models.Admin, error) {
	query := `SELECT id, username, password, full_name, is_active, created_at, updated_at 
              FROM admins WHERE username = $1`

	var admin models.Admin
	err := r.db.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Password, // Теперь читаем простой пароль
		&admin.FullName,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("администратор не найден")
		}
		return nil, fmt.Errorf("ошибка получения администратора: %w", err)
	}

	return &admin, nil
}

// GetAllTeachers получает всех учителей
func (r *AdminRepository) GetAllTeachers() ([]map[string]interface{}, error) {
	query := `SELECT id, email, full_name, is_active, created_at FROM teachers ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения учителей: %w", err)
	}
	defer rows.Close()

	var teachers []map[string]interface{}
	for rows.Next() {
		var id, email, fullName, isActive, createdAt interface{}
		err := rows.Scan(&id, &email, &fullName, &isActive, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		teachers = append(teachers, map[string]interface{}{
			"id":         id,
			"email":      email,
			"full_name":  fullName,
			"is_active":  isActive,
			"created_at": createdAt,
		})
	}

	return teachers, nil
}

// CreateTeacher создает нового учителя
func (r *AdminRepository) CreateTeacher(email, password, fullName string) error {
	query := `INSERT INTO teachers (email, password, full_name, is_active, created_at, updated_at) 
              VALUES ($1, $2, $3, TRUE, NOW(), NOW())`

	_, err := r.db.Exec(query, email, password, fullName)
	if err != nil {
		return fmt.Errorf("ошибка создания учителя: %w", err)
	}

	return nil
}

// UpdateTeacher обновляет учителя
func (r *AdminRepository) UpdateTeacher(id int, email, fullName string, isActive bool) error {
	query := `UPDATE teachers SET email = $1, full_name = $2, is_active = $3, updated_at = NOW() WHERE id = $4`

	result, err := r.db.Exec(query, email, fullName, isActive, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления учителя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("учитель не найден")
	}

	return nil
}

// DeleteTeacher удаляет учителя
func (r *AdminRepository) DeleteTeacher(id int) error {
	query := `DELETE FROM teachers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления учителя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("учитель не найден")
	}

	return nil
}

// GetTeacherByID получает учителя по ID
func (r *AdminRepository) GetTeacherByID(id int) (map[string]interface{}, error) {
	query := `SELECT id, email, full_name, is_active, created_at FROM teachers WHERE id = $1`

	var tID, email, fullName, isActive, createdAt interface{}
	err := r.db.QueryRow(query, id).Scan(&tID, &email, &fullName, &isActive, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("учитель не найден")
		}
		return nil, fmt.Errorf("ошибка получения учителя: %w", err)
	}

	return map[string]interface{}{
		"id":         tID,
		"email":      email,
		"full_name":  fullName,
		"is_active":  isActive,
		"created_at": createdAt,
	}, nil
}
