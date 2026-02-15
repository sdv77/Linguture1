package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sdv77/Linguture1/internal/models"
)

// UserLessonRepository работает с таблицей прогресса пользователей
type UserLessonRepository struct {
	db *sql.DB
}

// NewUserLessonRepository создает новый репозиторий
func NewUserLessonRepository(db *sql.DB) *UserLessonRepository {
	return &UserLessonRepository{db: db}
}

// GetOrCreate создает или получает прогресс пользователя по уроку
func (r *UserLessonRepository) GetOrCreate(userID, lessonID int) (*models.UserLesson, error) {
	// Пытаемся найти существующий прогресс
	query := `SELECT id, user_id, lesson_id, status, started_at, completed_at, score, created_at, updated_at 
              FROM user_lessons WHERE user_id = $1 AND lesson_id = $2`

	var userLesson models.UserLesson
	err := r.db.QueryRow(query, userID, lessonID).Scan(
		&userLesson.ID,
		&userLesson.UserID,
		&userLesson.LessonID,
		&userLesson.Status,
		&userLesson.StartedAt,
		&userLesson.CompletedAt,
		&userLesson.Score,
		&userLesson.CreatedAt,
		&userLesson.UpdatedAt,
	)

	if err == nil {
		// Прогресс найден
		return &userLesson, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("ошибка получения прогресса: %w", err)
	}

	// Создаем новый прогресс
	insertQuery := `INSERT INTO user_lessons (user_id, lesson_id, status, created_at, updated_at) 
                    VALUES ($1, $2, 'not_started', $3, $4) 
                    RETURNING id, user_id, lesson_id, status, started_at, completed_at, score, created_at, updated_at`

	now := time.Now()
	err = r.db.QueryRow(insertQuery, userID, lessonID, now, now).Scan(
		&userLesson.ID,
		&userLesson.UserID,
		&userLesson.LessonID,
		&userLesson.Status,
		&userLesson.StartedAt,
		&userLesson.CompletedAt,
		&userLesson.Score,
		&userLesson.CreatedAt,
		&userLesson.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания прогресса: %w", err)
	}

	return &userLesson, nil
}

// Update обновляет прогресс пользователя
func (r *UserLessonRepository) Update(id int, input models.UserLessonInput) error {
	// Получаем текущий прогресс для проверки статуса
	var currentStatus string
	err := r.db.QueryRow(`SELECT status FROM user_lessons WHERE id = $1`, id).Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("ошибка получения текущего статуса: %w", err)
	}

	now := time.Now()
	query := `UPDATE user_lessons SET status = $1, score = $2, updated_at = $3`

	// Если статус меняется на "in_progress" и раньше был "not_started", устанавливаем started_at
	if input.Status == "in_progress" && currentStatus == "not_started" {
		query += `, started_at = $4`
		query += ` WHERE id = $5`

		_, err = r.db.Exec(query, input.Status, input.Score, now, now, id)
	} else if input.Status == "completed" && currentStatus != "completed" {
		// Если статус меняется на "completed", устанавливаем completed_at
		query += `, completed_at = $4`
		query += ` WHERE id = $5`

		_, err = r.db.Exec(query, input.Status, input.Score, now, now, id)
	} else {
		query += ` WHERE id = $4`
		_, err = r.db.Exec(query, input.Status, input.Score, now, id)
	}

	if err != nil {
		return fmt.Errorf("ошибка обновления прогресса: %w", err)
	}

	return nil
}

// GetByUserID получает все уроки пользователя с прогрессом
func (r *UserLessonRepository) GetByUserID(userID int) ([]models.UserLesson, error) {
	query := `SELECT id, user_id, lesson_id, status, started_at, completed_at, score, created_at, updated_at 
              FROM user_lessons WHERE user_id = $1 ORDER BY lesson_id`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения прогресса пользователя: %w", err)
	}
	defer rows.Close()

	var userLessons []models.UserLesson
	for rows.Next() {
		var userLesson models.UserLesson
		err := rows.Scan(
			&userLesson.ID,
			&userLesson.UserID,
			&userLesson.LessonID,
			&userLesson.Status,
			&userLesson.StartedAt,
			&userLesson.CompletedAt,
			&userLesson.Score,
			&userLesson.CreatedAt,
			&userLesson.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		userLessons = append(userLessons, userLesson)
	}

	return userLessons, nil
}

// GetProgressStats получает статистику прогресса пользователя
func (r *UserLessonRepository) GetProgressStats(userID int) (*models.UserLessonProgress, error) {
	query := `
        SELECT 
            COUNT(*) as total,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
            COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress,
            COUNT(CASE WHEN status = 'not_started' THEN 1 END) as not_started,
            COALESCE(SUM(score), 0) as total_score
        FROM user_lessons 
        WHERE user_id = $1
    `

	var stats models.UserLessonProgress
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalLessons,
		&stats.CompletedLessons,
		&stats.InProgressLessons,
		&stats.NotStartedLessons,
		&stats.TotalScore,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения статистики: %w", err)
	}

	return &stats, nil
}
