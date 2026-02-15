package service

import (
	"time"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
)

// WordService содержит бизнес-логику
type WordService struct {
	repo *repository.WordRepository
}

// NewWordService создает новый сервис
func NewWordService(repo *repository.WordRepository) *WordService {
	return &WordService{repo: repo}
}

// CreateWord создает новое слово для пользователя
func (s *WordService) CreateWord(userID int, input models.WordInput) (*models.Word, error) {
	word := &models.Word{
		UserID:    userID,
		Word:      input.Word,
		Meaning:   input.Meaning,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(word); err != nil {
		return nil, err
	}

	return word, nil
}

// GetAllWords получает все слова пользователя
func (s *WordService) GetAllWords(userID int) ([]models.Word, error) {
	return s.repo.GetAll(userID)
}

// GetWordByID получает слово по ID
func (s *WordService) GetWordByID(id int, userID int) (*models.Word, error) {
	return s.repo.GetByID(id, userID)
}

// UpdateWord обновляет слово пользователя
func (s *WordService) UpdateWord(id int, userID int, input models.WordInput) error {
	word := &models.Word{
		Word:      input.Word,
		Meaning:   input.Meaning,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(id, userID, word)
}

// DeleteWord удаляет слово пользователя
func (s *WordService) DeleteWord(id int, userID int) error {
	return s.repo.Delete(id, userID)
}
