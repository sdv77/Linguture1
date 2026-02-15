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

// CreateWord создает новое слово
func (s *WordService) CreateWord(input models.WordInput) (*models.Word, error) {
	word := &models.Word{
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

// GetAllWords получает все слова
func (s *WordService) GetAllWords() ([]models.Word, error) {
	return s.repo.GetAll()
}

// GetWordByID получает слово по ID
func (s *WordService) GetWordByID(id int) (*models.Word, error) {
	return s.repo.GetByID(id)
}

// UpdateWord обновляет слово
func (s *WordService) UpdateWord(id int, input models.WordInput) error {
	word := &models.Word{
		Word:      input.Word,
		Meaning:   input.Meaning,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(id, word)
}

// DeleteWord удаляет слово
func (s *WordService) DeleteWord(id int) error {
	return s.repo.Delete(id)
}
