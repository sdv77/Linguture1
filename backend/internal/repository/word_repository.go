package repository

import (
	"database/sql"
	"fmt"

	"github.com/sdv77/Linguture1/internal/models"
)

// WordRepository работает с базой данных
type WordRepository struct {
	db *sql.DB
}

// NewWordRepository создает новый репозиторий
func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

// Create создает новое слово в базе
func (r *WordRepository) Create(word *models.Word) error {
	query := `INSERT INTO words (word, meaning, created_at, updated_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(query, word.Word, word.Meaning, word.CreatedAt, word.UpdatedAt).
		Scan(&word.ID)

	if err != nil {
		return fmt.Errorf("ошибка создания слова: %w", err)
	}

	return nil
}

// GetAll получает все слова
func (r *WordRepository) GetAll() ([]models.Word, error) {
	query := `SELECT id, word, meaning, created_at, updated_at FROM words ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слов: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.Word, &word.Meaning, &word.CreatedAt, &word.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

// GetByID получает слово по ID
func (r *WordRepository) GetByID(id int) (*models.Word, error) {
	query := `SELECT id, word, meaning, created_at, updated_at FROM words WHERE id = $1`

	var word models.Word
	err := r.db.QueryRow(query, id).Scan(&word.ID, &word.Word, &word.Meaning, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("слово не найдено")
		}
		return nil, fmt.Errorf("ошибка получения слова: %w", err)
	}

	return &word, nil
}

// Update обновляет слово
func (r *WordRepository) Update(id int, word *models.Word) error {
	query := `UPDATE words SET word = $1, meaning = $2, updated_at = $3 WHERE id = $4`

	result, err := r.db.Exec(query, word.Word, word.Meaning, word.UpdatedAt, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления слова: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("слово не найдено")
	}

	return nil
}

// Delete удаляет слово
func (r *WordRepository) Delete(id int) error {
	query := `DELETE FROM words WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления слова: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("слово не найдено")
	}

	return nil
}
