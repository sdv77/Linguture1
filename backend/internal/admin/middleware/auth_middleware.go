package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sdv77/Linguture1/pkg/token"
)

type AdminAuthMiddleware struct {
	tokenService *token.Service
}

func NewAdminAuthMiddleware(tokenService *token.Service) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{tokenService: tokenService}
}

func (m *AdminAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		_, _, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Недействительный токен: %v", err), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
