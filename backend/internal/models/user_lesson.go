package models

import (
	"time"
)

// UserLesson представляет прогресс пользователя по уроку
type UserLesson struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	LessonID    int        `json:"lesson_id"`
	Status      string     `json:"status"` // not_started, in_progress, completed
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Score       int        `json:"score"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserLessonInput данные для обновления прогресса
type UserLessonInput struct {
	Status string `json:"status"`
	Score  int    `json:"score"`
}

// UserLessonProgress статистика прогресса пользователя
type UserLessonProgress struct {
	TotalLessons      int `json:"total_lessons"`
	CompletedLessons  int `json:"completed_lessons"`
	InProgressLessons int `json:"in_progress_lessons"`
	NotStartedLessons int `json:"not_started_lessons"`
	TotalScore        int `json:"total_score"`
}
