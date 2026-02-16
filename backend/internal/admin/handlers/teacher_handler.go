package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sdv77/Linguture1/internal/admin/service"

	"github.com/gorilla/mux"
)

type TeacherHandler struct {
	adminService *service.AdminService
}

func NewTeacherHandler(adminService *service.AdminService) *TeacherHandler {
	return &TeacherHandler{adminService: adminService}
}

func (h *TeacherHandler) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.adminService.GetAllTeachers()
	if err != nil {
		log.Printf("Ошибка получения учителей: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// CreateTeacher обрабатывает создание учителя
func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"` // ИСПРАВЛЕНО: было "full_name" в JSON
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Если полное имя не указано, используем email
	fullName := input.FullName
	if fullName == "" {
		fullName = input.Email
	}

	err := h.adminService.CreateTeacher(input.Email, input.Password, fullName)
	if err != nil {
		log.Printf("Ошибка создания учителя: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Учитель успешно создан",
	})
}

// UpdateTeacher обрабатывает обновление учителя
func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"` // ИСПРАВЛЕНО
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Если полное имя пустое, оставляем как есть
	fullName := input.FullName
	if fullName == "" {
		fullName = input.Email
	}

	err = h.adminService.UpdateTeacher(id, input.Email, fullName, input.IsActive)
	if err != nil {
		log.Printf("Ошибка обновления учителя: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Учитель успешно обновлен",
	})
}

func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.adminService.DeleteTeacher(id)
	if err != nil {
		log.Printf("Ошибка удаления учителя: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Учитель успешно удален",
	})
}

// GetTeacherByID обрабатывает получение учителя по ID
func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.adminService.GetTeacherByID(id)
	if err != nil {
		log.Printf("Ошибка получения учителя: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}
