package models

import (
	"time"
)

// Lesson представляет урок
type Lesson struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       int       `json:"level"`
	LessonType  string    `json:"lesson_type"`
	IsActive    bool      `json:"is_active"`
	OrderNum    int       `json:"order_num"`
	TeacherID   int       `json:"teacher_id,omitempty"` // Кто создал урок
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LessonVocabulary представляет слово из урока
type LessonVocabulary struct {
	ID            int       `json:"id"`
	Word          string    `json:"word"`
	Meaning       string    `json:"meaning"`
	Transcription string    `json:"transcription,omitempty"`
	Example       string    `json:"example,omitempty"`
	LessonID      int       `json:"lesson_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// LessonWithWords представляет урок со словами
type LessonWithWords struct {
	Lesson Lesson             `json:"lesson"`
	Words  []LessonVocabulary `json:"words"`
}

// LessonInput данные для создания/обновления урока
type LessonInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	LessonType  string `json:"lesson_type"`
	OrderNum    int    `json:"order_num"`
}

// LessonVocabularyInput данные для слова урока
type LessonVocabularyInput struct {
	Word          string `json:"word" validate:"required"`
	Meaning       string `json:"meaning" validate:"required"`
	Transcription string `json:"transcription"`
	Example       string `json:"example"`
	LessonID      int    `json:"lesson_id" validate:"required"`
}
