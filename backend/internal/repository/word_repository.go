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
	query := `INSERT INTO words (user_id, word, meaning, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(query, word.UserID, word.Word, word.Meaning, word.CreatedAt, word.UpdatedAt).
		Scan(&word.ID)

	if err != nil {
		return fmt.Errorf("ошибка создания слова: %w", err)
	}

	return nil
}

// GetAll получает все слова пользователя
func (r *WordRepository) GetAll(userID int) ([]models.Word, error) {
	query := `SELECT id, user_id, word, meaning, created_at, updated_at 
              FROM words WHERE user_id = $1 ORDER BY id`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слов: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.UserID, &word.Word, &word.Meaning, &word.CreatedAt, &word.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

// GetByID получает слово по ID и проверяет принадлежность пользователю
func (r *WordRepository) GetByID(id int, userID int) (*models.Word, error) {
	query := `SELECT id, user_id, word, meaning, created_at, updated_at 
              FROM words WHERE id = $1 AND user_id = $2`

	var word models.Word
	err := r.db.QueryRow(query, id, userID).Scan(&word.ID, &word.UserID, &word.Word, &word.Meaning, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("слово не найдено")
		}
		return nil, fmt.Errorf("ошибка получения слова: %w", err)
	}

	return &word, nil
}

// Update обновляет слово пользователя
func (r *WordRepository) Update(id int, userID int, word *models.Word) error {
	query := `UPDATE words SET word = $1, meaning = $2, updated_at = $3 
              WHERE id = $4 AND user_id = $5`

	result, err := r.db.Exec(query, word.Word, word.Meaning, word.UpdatedAt, id, userID)
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

// Delete удаляет слово пользователя
func (r *WordRepository) Delete(id int, userID int) error {
	query := `DELETE FROM words WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
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
