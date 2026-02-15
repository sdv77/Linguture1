package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/service"
	"github.com/sdv77/Linguture1/pkg/email"
)

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authService  *service.AuthService
	emailService *email.Service
}

// NewAuthHandler создает новый обработчик
func NewAuthHandler(authService *service.AuthService, emailService *email.Service) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
	}
}

// Register обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.UserRegisterInput

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Регистрируем пользователя
	if err := h.authService.Register(input); err != nil {
		if err.Error() == "пользователь с таким email уже существует" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Printf("Ошибка регистрации: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован. Проверьте почту для подтверждения.",
	})
}

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.UserLoginInput

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Выполняем вход
	tokenResp, err := h.authService.Login(input)
	if err != nil {
		log.Printf("Ошибка входа: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}

// VerifyEmail обрабатывает подтверждение почты
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из query параметра
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Токен не предоставлен", http.StatusBadRequest)
		return
	}

	// Подтверждаем почту
	if err := h.authService.VerifyEmail(token); err != nil {
		log.Printf("Ошибка подтверждения почты: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Почта успешно подтверждена! Теперь вы можете войти в систему.",
	})
}

// GetCurrentUser обрабатывает получение данных текущего пользователя
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из заголовка (добавленного в middleware)
	// В реальном приложении используем контекст
	// Здесь для простоты передадим фиктивный ID
	// TODO: Получить реальный userID из токена

	// Для примера используем ID из тестового пользователя
	userID := 1

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
