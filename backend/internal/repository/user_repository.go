package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sdv77/Linguture1/internal/models"
)

// UserRepository работает с таблицей пользователей
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый репозиторий
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password_hash, is_verified, verification_token, verification_token_expires, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(query,
		user.Email,
		user.PasswordHash,
		true,
		user.VerificationToken,
		user.VerificationTokenExpires,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return nil
}

// FindByEmail находит пользователя по email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, is_verified, verification_token, verification_token_expires, created_at, updated_at 
              FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.VerificationToken,
		&user.VerificationTokenExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return &user, nil
}

// FindByID находит пользователя по ID
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, email, password_hash, is_verified, verification_token, verification_token_expires, created_at, updated_at 
              FROM users WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.VerificationToken,
		&user.VerificationTokenExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return &user, nil
}

// VerifyEmail подтверждает email пользователя
func (r *UserRepository) VerifyEmail(token string) error {
	query := `UPDATE users 
              SET is_verified = TRUE, verification_token = NULL, verification_token_expires = NULL, updated_at = $2 
              WHERE verification_token = $1 AND verification_token_expires > $2 AND is_verified = FALSE`

	now := time.Now()
	result, err := r.db.Exec(query, token, now)
	if err != nil {
		return fmt.Errorf("ошибка подтверждения email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("недействительный или просроченный токен")
	}

	return nil
}

// UpdateVerificationToken обновляет токен подтверждения
func (r *UserRepository) UpdateVerificationToken(userID int, token string, expiresAt time.Time) error {
	query := `UPDATE users SET verification_token = $1, verification_token_expires = $2, updated_at = $3 WHERE id = $4`

	_, err := r.db.Exec(query, token, expiresAt, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления токена подтверждения: %w", err)
	}

	return nil
}
