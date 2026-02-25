package models

import (
	"time"
)

// User представляет пользователя
type User struct {
	ID                       int        `json:"id"`
	Email                    string     `json:"email"`
	PasswordHash             string     `json:"-"` // Не отправляем хэш в ответе — это безопасность
	IsVerified               bool       `json:"is_verified"`
	VerificationToken        string     `json:"-"`
	VerificationTokenExpires *time.Time `json:"-"`

	// 🔹 НОВЫЕ ПОЛЯ для настройки профиля 🔹
	Nickname         string `json:"nickname,omitempty"`          // Уникальный никнейм пользователя (показываем в профиле)
	NativeLanguage   string `json:"native_language,omitempty"`   // Код языка (например "ru", "en") — родной язык
	LearningLanguage string `json:"learning_language,omitempty"` // Код языка, который пользователь учит
	IsSetup          bool   `json:"is_setup"`                    // ✅ true = профиль настроен, false = нужно донастроить

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRegisterInput данные для регистрации (без изменений)
type UserRegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginInput данные для входа (без изменений)
type UserLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// 🔹 НОВАЯ структура: данные для первичной настройки профиля 🔹
type UserProfileSetupInput struct {
	Nickname         string `json:"nickname" validate:"required,min=3,max=30"`
	NativeLanguage   string `json:"native_language" validate:"required,len=2"`   // например "ru", "en"
	LearningLanguage string `json:"learning_language" validate:"required,len=2"` // например "en", "es"
}

// UserResponse данные пользователя для ответа (обновляем)
type UserResponse struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Nickname         string    `json:"nickname,omitempty"`
	NativeLanguage   string    `json:"native_language,omitempty"`
	LearningLanguage string    `json:"learning_language,omitempty"`
	IsVerified       bool      `json:"is_verified"`
	IsSetup          bool      `json:"is_setup"` // 🔹 Важно: фронтенд будет проверять это поле
	CreatedAt        time.Time `json:"created_at"`
}

// TokenResponse ответ с токеном (без изменений)
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}
