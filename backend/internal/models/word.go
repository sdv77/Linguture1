package models

import "time"

// Word представляет сущность "слово"
type Word struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // ID пользователя
	Word      string    `json:"word"`
	Meaning   string    `json:"meaning"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WordInput используется для создания/обновления слова
type WordInput struct {
	Word    string `json:"word" validate:"required"`
	Meaning string `json:"meaning" validate:"required"`
}
