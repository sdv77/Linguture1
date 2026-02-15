package models

import "time"

// Word представляет сущность "слово"
type Word struct {
	ID        int       `json:"id"`         // Уникальный идентификатор
	Word      string    `json:"word"`       // Само слово
	Meaning   string    `json:"meaning"`    // Значение слова
	CreatedAt time.Time `json:"created_at"` // Дата создания
	UpdatedAt time.Time `json:"updated_at"` // Дата обновления
}

// WordInput используется для создания/обновления слова
type WordInput struct {
	Word    string `json:"word" validate:"required"`
	Meaning string `json:"meaning" validate:"required"`
}
