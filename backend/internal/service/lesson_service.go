package service

import (
	"time"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
)

// LessonService содержит бизнес-логику уроков
type LessonService struct {
	lessonRepo *repository.LessonRepository
}

// NewLessonService создает новый сервис
func NewLessonService(lessonRepo *repository.LessonRepository) *LessonService {
	return &LessonService{lessonRepo: lessonRepo}
}

// GetAllLessons получает все уроки
func (s *LessonService) GetAllLessons() ([]models.Lesson, error) {
	return s.lessonRepo.GetAll()
}

// GetLessonByID получает урок по ID
func (s *LessonService) GetLessonByID(id int) (*models.Lesson, error) {
	return s.lessonRepo.GetByID(id)
}

// GetLessonWithWords получает урок со словами
func (s *LessonService) GetLessonWithWords(lessonID int) (*models.LessonWithWords, error) {
	return s.lessonRepo.GetLessonWithWords(lessonID)
}

// CreateLesson создает новый урок
func (s *LessonService) CreateLesson(input models.LessonInput, teacherID int) (*models.Lesson, error) {
	lesson := &models.Lesson{
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		LessonType:  input.LessonType,
		OrderNum:    input.OrderNum,
		TeacherID:   teacherID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.lessonRepo.Create(lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

// UpdateLesson обновляет урок
func (s *LessonService) UpdateLesson(id int, input models.LessonInput, teacherID int) error {
	lesson := &models.Lesson{
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		LessonType:  input.LessonType,
		OrderNum:    input.OrderNum,
		TeacherID:   teacherID,
		UpdatedAt:   time.Now(),
	}

	return s.lessonRepo.Update(id, lesson)
}

// DeleteLesson удаляет урок
func (s *LessonService) DeleteLesson(id int) error {
	return s.lessonRepo.Delete(id)
}

// GetTeacherLessons получает уроки учителя
func (s *LessonService) GetTeacherLessons(teacherID int) ([]models.Lesson, error) {
	return s.lessonRepo.GetByTeacherID(teacherID)
}
