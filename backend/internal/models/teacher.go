package models

import (
	"time"
)

// Teacher представляет учителя
type Teacher struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Пароль (не отправляем в ответе)
	PasswordHash string    `json:"-"` // Хэш пароля (для совместимости)
	FullName     string    `json:"full_name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TeacherLoginInput данные для входа учителя
type TeacherLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TeacherResponse данные учителя для ответа
type TeacherResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
