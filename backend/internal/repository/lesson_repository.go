package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sdv77/Linguture1/internal/models"
)

// LessonRepository работает с таблицей уроков
type LessonRepository struct {
	db *sql.DB
}

// NewLessonRepository создает новый репозиторий
func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// GetAll получает все активные уроки
func (r *LessonRepository) GetAll() ([]models.Lesson, error) {
	query := `SELECT id, title, description, level, lesson_type, is_active, order_num, created_at, updated_at 
              FROM lessons 
              WHERE is_active = TRUE 
              ORDER BY order_num, level`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения уроков: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Level,
			&lesson.LessonType,
			&lesson.IsActive,
			&lesson.OrderNum,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// GetByID получает урок по ID
func (r *LessonRepository) GetByID(id int) (*models.Lesson, error) {
	query := `SELECT id, title, description, level, lesson_type, is_active, order_num, created_at, updated_at 
              FROM lessons WHERE id = $1`

	var lesson models.Lesson
	err := r.db.QueryRow(query, id).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.Description,
		&lesson.Level,
		&lesson.LessonType,
		&lesson.IsActive,
		&lesson.OrderNum,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("урок не найден")
		}
		return nil, fmt.Errorf("ошибка получения урока: %w", err)
	}

	return &lesson, nil
}

// GetWordsByLessonID получает слова урока по ID урока
func (r *LessonRepository) GetWordsByLessonID(lessonID int) ([]models.LessonVocabulary, error) {
	query := `SELECT id, word, meaning, transcription, example, lesson_id, created_at 
              FROM lesson_vocabulary 
              WHERE lesson_id = $1 
              ORDER BY order_in_lesson, id`

	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слов урока: %w", err)
	}
	defer rows.Close()

	var words []models.LessonVocabulary
	for rows.Next() {
		var word models.LessonVocabulary
		err := rows.Scan(
			&word.ID,
			&word.Word,
			&word.Meaning,
			&word.Transcription,
			&word.Example,
			&word.LessonID,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

// GetLessonWithWords получает урок со всеми словами
func (r *LessonRepository) GetLessonWithWords(lessonID int) (*models.LessonWithWords, error) {
	lesson, err := r.GetByID(lessonID)
	if err != nil {
		return nil, err
	}

	words, err := r.GetWordsByLessonID(lessonID)
	if err != nil {
		return nil, err
	}

	return &models.LessonWithWords{
		Lesson: *lesson,
		Words:  words,
	}, nil
}

// Create создает новый урок
func (r *LessonRepository) Create(lesson *models.Lesson) error {
	query := `INSERT INTO lessons (title, description, level, lesson_type, is_active, order_num, teacher_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := r.db.QueryRow(query,
		lesson.Title,
		lesson.Description,
		lesson.Level,
		lesson.LessonType,
		lesson.IsActive,
		lesson.OrderNum,
		lesson.TeacherID,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	).Scan(&lesson.ID)

	if err != nil {
		return fmt.Errorf("ошибка создания урока: %w", err)
	}

	return nil
}

// Update обновляет урок
func (r *LessonRepository) Update(id int, lesson *models.Lesson) error {
	query := `UPDATE lessons SET title = $1, description = $2, level = $3, lesson_type = $4, 
              is_active = $5, order_num = $6, teacher_id = $7, updated_at = $8 
              WHERE id = $9`

	result, err := r.db.Exec(query,
		lesson.Title,
		lesson.Description,
		lesson.Level,
		lesson.LessonType,
		lesson.IsActive,
		lesson.OrderNum,
		lesson.TeacherID,
		lesson.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("ошибка обновления урока: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("урок не найден")
	}

	return nil
}

// Delete удаляет урок
func (r *LessonRepository) Delete(id int) error {
	query := `DELETE FROM lessons WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления урока: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("урок не найден")
	}

	return nil
}

// GetByTeacherID получает уроки учителя
func (r *LessonRepository) GetByTeacherID(teacherID int) ([]models.Lesson, error) {
	query := `SELECT id, title, description, level, lesson_type, is_active, order_num, teacher_id, created_at, updated_at 
              FROM lessons WHERE teacher_id = $1 ORDER BY order_num, level`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения уроков учителя: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Level,
			&lesson.LessonType,
			&lesson.IsActive,
			&lesson.OrderNum,
			&lesson.TeacherID,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// AddWordToLesson добавляет слово к уроку
func (r *LessonRepository) AddWordToLesson(lessonID int, wordInput models.LessonVocabularyInput) error {
	query := `INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id, order_in_lesson, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		wordInput.Word,
		wordInput.Meaning,
		wordInput.Transcription,
		wordInput.Example,
		lessonID,
		wordInput.OrderInLesson,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("ошибка добавления слова к уроку: %w", err)
	}

	return nil
}

// DeleteWordFromLesson удаляет слово из урока
func (r *LessonRepository) DeleteWordFromLesson(wordID int) error {
	query := `DELETE FROM lesson_vocabulary WHERE id = $1`

	result, err := r.db.Exec(query, wordID)
	if err != nil {
		return fmt.Errorf("ошибка удаления слова из урока: %w", err)
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
