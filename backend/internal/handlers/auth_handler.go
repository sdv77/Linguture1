package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.UserLoginInput

	// 🔹 1. Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// 🔹 2. Простая валидация: email и пароль не пустые
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// 🔹 3. Выполняем вход через сервис (проверка пароля, генерация токена)
	tokenResp, err := h.authService.Login(input)
	if err != nil {
		log.Printf("Ошибка входа: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 🔹 4. Получаем данные пользователя для ответа
	// Поскольку Login — публичный маршрут (без middleware),
	// мы не можем взять userID из заголовка.
	// Вместо этого запрашиваем пользователя по email (который мы уже проверили)
	user, err := h.authService.GetUserByEmail(input.Email)
	if err != nil {
		// 🔹 Важно: не прерываем вход, если не удалось получить доп. данные
		// Возвращаем хотя бы токен, чтобы пользователь мог войти
		log.Printf("Предупреждение: не удалось получить данные пользователя после входа: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResp)
		return
	}

	// 🔹 5. Формируем UserResponse с новыми полями профиля
	userResp := &models.UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		Nickname:         user.Nickname,         // 🔹 новое поле
		NativeLanguage:   user.NativeLanguage,   // 🔹 новое поле
		LearningLanguage: user.LearningLanguage, // 🔹 новое поле
		IsVerified:       user.IsVerified,
		IsSetup:          user.IsSetup, // 🔹 критически важно для фронтенда!
		CreatedAt:        user.CreatedAt,
	}

	// 🔹 6. Формируем итоговый ответ
	// Флаг needs_setup — удобный булеан для фронтенда:
	// если true → сразу показать окно настройки профиля
	response := map[string]interface{}{
		"token":       tokenResp.Token,
		"expires_in":  tokenResp.ExpiresIn,
		"user":        userResp,
		"needs_setup": !user.IsSetup, // 🔹 инвертируем: если IsSetup=false → needs_setup=true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

// SetupProfile обрабатывает первичную настройку профиля
func (h *AuthHandler) SetupProfile(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из заголовка (добавленного middleware)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	var input models.UserProfileSetupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация на уровне хендлера
	if input.Nickname == "" || input.NativeLanguage == "" || input.LearningLanguage == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	// Вызываем сервис
	if err := h.authService.SetupProfile(userID, input); err != nil {
		log.Printf("Ошибка настройки профиля: %v", err)

		// Разные ошибки → разные HTTP-коды
		if err.Error() == "никнейм уже занят" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if err.Error() == "пользователь не найден или уже настроен" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Профиль успешно настроен!",
		"is_setup": true,
	})
}
