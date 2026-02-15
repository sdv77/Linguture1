package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/service"
)

// TeacherAuthHandler обрабатывает запросы аутентификации учителей
type TeacherAuthHandler struct {
	authService *service.TeacherAuthService
}

// NewTeacherAuthHandler создает новый обработчик
func NewTeacherAuthHandler(authService *service.TeacherAuthService) *TeacherAuthHandler {
	return &TeacherAuthHandler{authService: authService}
}

// Login обрабатывает вход учителя
func (h *TeacherAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.TeacherLoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.authService.Login(input)
	if err != nil {
		log.Printf("Ошибка входа учителя: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}

// GetTeacher обрабатывает получение данных текущего учителя
func (h *TeacherAuthHandler) GetTeacher(w http.ResponseWriter, r *http.Request) {
	teacherID := 1 // TODO: Получить реальный teacherID из токена

	teacher, err := h.authService.GetTeacherByID(teacherID)
	if err != nil {
		log.Printf("Ошибка получения учителя: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}
