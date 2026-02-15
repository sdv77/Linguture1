package service

import (
	"fmt"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
)

// UserLessonService содержит бизнес-логику прогресса пользователя
type UserLessonService struct {
	userLessonRepo *repository.UserLessonRepository
	lessonRepo     *repository.LessonRepository
}

// NewUserLessonService создает новый сервис
func NewUserLessonService(userLessonRepo *repository.UserLessonRepository, lessonRepo *repository.LessonRepository) *UserLessonService {
	return &UserLessonService{
		userLessonRepo: userLessonRepo,
		lessonRepo:     lessonRepo,
	}
}

// StartLesson начинает урок для пользователя
func (s *UserLessonService) StartLesson(userID, lessonID int) error {
	// Проверяем, существует ли урок
	_, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return fmt.Errorf("урок не найден: %w", err)
	}

	// Получаем или создаём прогресс
	userLesson, err := s.userLessonRepo.GetOrCreate(userID, lessonID)
	if err != nil {
		return err
	}

	// Если урок уже завершён, не меняем статус
	if userLesson.Status == "completed" {
		return nil
	}

	// Обновляем статус на "in_progress"
	input := models.UserLessonInput{
		Status: "in_progress",
		Score:  userLesson.Score,
	}

	return s.userLessonRepo.Update(userLesson.ID, input)
}

// CompleteLesson завершает урок для пользователя
func (s *UserLessonService) CompleteLesson(userID, lessonID, score int) error {
	// Проверяем, существует ли урок
	_, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return fmt.Errorf("урок не найден: %w", err)
	}

	// Получаем прогресс
	userLesson, err := s.userLessonRepo.GetOrCreate(userID, lessonID)
	if err != nil {
		return err
	}

	// Обновляем статус на "completed"
	input := models.UserLessonInput{
		Status: "completed",
		Score:  score,
	}

	return s.userLessonRepo.Update(userLesson.ID, input)
}

// GetUserLessons получает все уроки пользователя с прогрессом
func (s *UserLessonService) GetUserLessons(userID int) ([]models.UserLesson, error) {
	return s.userLessonRepo.GetByUserID(userID)
}

// GetProgressStats получает статистику прогресса пользователя
func (s *UserLessonService) GetProgressStats(userID int) (*models.UserLessonProgress, error) {
	return s.userLessonRepo.GetProgressStats(userID)
}
