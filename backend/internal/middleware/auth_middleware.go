package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sdv77/Linguture1/pkg/token"
)

// AuthMiddleware проверяет наличие и валидность JWT токена
type AuthMiddleware struct {
	tokenService *token.Service
}

// NewAuthMiddleware создает новый middleware
func NewAuthMiddleware(tokenService *token.Service) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

// Middleware возвращает функцию-обработчик для проверки токена

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 🔹 ВАЖНО: пропускаем OPTIONS запросы (CORS preflight)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Проверяем токен
		userID, email, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Недействительный токен: %v", err), http.StatusUnauthorized)
			return
		}

		// Добавляем данные пользователя в заголовок
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		r.Header.Set("X-User-Email", email)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// GetUserID извлекает ID пользователя из заголовка запроса
func GetUserID(r *http.Request) int {
	// В реальном приложении лучше использовать контекст
	// Это упрощённая версия
	return 0 // Будем передавать userID напрямую в хендлерах
}
