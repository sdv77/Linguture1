package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service управляет JWT токенами
type Service struct {
	secretKey               string
	tokenExpiry             time.Duration
	tokenExpiryVerification time.Duration
}

// NewService создает новый сервис токенов
func NewService(secretKey string, tokenExpiry time.Duration) *Service {
	return &Service{
		secretKey:               secretKey,
		tokenExpiry:             tokenExpiry,
		tokenExpiryVerification: 24 * time.Hour, // Токен подтверждения действителен 24 часа
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (s *Service) GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка создания токена: %w", err)
	}

	return tokenString, nil
}

// VerifyToken проверяет и парсит JWT токен
func (s *Service) VerifyToken(tokenStr string) (int, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("недопустимый метод подписи")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("недействительный токен: %w", err)
	}

	if !token.Valid {
		return 0, "", errors.New("токен недействителен")
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("невозможно получить данные из токена")
	}

	// Получаем данные
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("невозможно получить user_id из токена")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", errors.New("невозможно получить email из токена")
	}

	// Конвертируем в целое число
	userID := int(userIDFloat)

	return userID, email, nil
}

// GenerateVerificationToken генерирует токен для подтверждения почты
func (s *Service) GenerateVerificationToken() (string, error) {
	// Используем случайную строку как токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(s.tokenExpiryVerification).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка создания токена подтверждения: %w", err)
	}

	return tokenString, nil
}

// GetTokenExpiry возвращает время жизни токена
func (s *Service) GetTokenExpiry() time.Duration {
	return s.tokenExpiry
}
