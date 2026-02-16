package service

import (
	"fmt"
	"log"
	"time"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
)

// UserLessonService содержит бизнес-логику прогресса пользователя
type UserLessonService struct {
	userLessonRepo *repository.UserLessonRepository
	lessonRepo     *repository.LessonRepository
	wordRepo       *repository.WordRepository // Добавили
}

// NewUserLessonService создает новый сервис
func NewUserLessonService(userLessonRepo *repository.UserLessonRepository, lessonRepo *repository.LessonRepository, wordRepo *repository.WordRepository) *UserLessonService {
	return &UserLessonService{
		userLessonRepo: userLessonRepo,
		lessonRepo:     lessonRepo,
		wordRepo:       wordRepo,
	}
}

// StartLesson начинает урок для пользователя
func (s *UserLessonService) StartLesson(userID, lessonID int) error {
	_, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return fmt.Errorf("урок не найден: %w", err)
	}

	userLesson, err := s.userLessonRepo.GetOrCreate(userID, lessonID)
	if err != nil {
		return err
	}

	if userLesson.Status == "completed" {
		return nil
	}

	input := models.UserLessonInput{
		Status: "in_progress",
		Score:  userLesson.Score,
	}

	return s.userLessonRepo.Update(userLesson.ID, input)
}

// CompleteLesson завершает урок для пользователя и добавляет слова в профиль
func (s *UserLessonService) CompleteLesson(userID, lessonID, score int) error {
	_, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return fmt.Errorf("урок не найден: %w", err)
	}

	userLesson, err := s.userLessonRepo.GetOrCreate(userID, lessonID)
	if err != nil {
		return err
	}

	input := models.UserLessonInput{
		Status: "completed",
		Score:  score,
	}

	if err := s.userLessonRepo.Update(userLesson.ID, input); err != nil {
		return err
	}

	// Получаем слова урока
	words, err := s.lessonRepo.GetWordsByLessonID(lessonID)
	if err != nil {
		log.Printf("Предупреждение: не удалось получить слова урока: %v", err)
		return nil
	}

	// Добавляем слова в профиль пользователя
	addedCount := 0
	for _, word := range words {
		// Проверяем, не существует ли уже такое слово у пользователя
		exists, err := s.wordRepo.ExistsByWordAndUserID(word.Word, userID)
		if err != nil {
			log.Printf("Ошибка проверки существования слова '%s': %v", word.Word, err)
			continue
		}

		if exists {
			log.Printf("Слово '%s' уже существует в профиле пользователя %d", word.Word, userID)
			continue
		}

		// Создаем новое слово для пользователя
		userWord := &models.Word{
			UserID:    userID,
			Word:      word.Word,
			Meaning:   word.Meaning,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.wordRepo.Create(userWord); err != nil {
			log.Printf("Ошибка добавления слова '%s' в профиль: %v", word.Word, err)
			continue
		}

		addedCount++
		log.Printf("Слово '%s' добавлено в профиль пользователя %d", word.Word, userID)
	}

	log.Printf("Добавлено %d новых слов в профиль пользователя %d", addedCount, userID)
	return nil
}

// GetUserLessons получает все уроки пользователя с прогрессом
func (s *UserLessonService) GetUserLessons(userID int) ([]models.UserLesson, error) {
	return s.userLessonRepo.GetByUserID(userID)
}

// GetProgressStats получает статистику прогресса пользователя
func (s *UserLessonService) GetProgressStats(userID int) (*models.UserLessonProgress, error) {
	return s.userLessonRepo.GetProgressStats(userID)
}
