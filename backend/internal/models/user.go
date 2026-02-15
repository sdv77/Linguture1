package models

import (
	"time"
)

// User представляет пользователя
type User struct {
	ID                       int        `json:"id"`
	Email                    string     `json:"email"`
	PasswordHash             string     `json:"-"` // Не отправляем хэш в ответе
	IsVerified               bool       `json:"is_verified"`
	VerificationToken        string     `json:"-"`
	VerificationTokenExpires *time.Time `json:"-"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

// UserRegisterInput данные для регистрации
type UserRegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginInput данные для входа
type UserLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse данные пользователя для ответа
type UserResponse struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// TokenResponse ответ с токеном
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}
