package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/service"
)

// LessonHandler обрабатывает HTTP запросы уроков
type LessonHandler struct {
	lessonService     *service.LessonService
	userLessonService *service.UserLessonService
}

// NewLessonHandler создает новый обработчик
func NewLessonHandler(lessonService *service.LessonService, userLessonService *service.UserLessonService) *LessonHandler {
	return &LessonHandler{
		lessonService:     lessonService,
		userLessonService: userLessonService,
	}
}

// GetAllLessons обрабатывает получение всех уроков
func (h *LessonHandler) GetAllLessons(w http.ResponseWriter, r *http.Request) {
	lessons, err := h.lessonService.GetAllLessons()
	if err != nil {
		log.Printf("Ошибка получения уроков: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons)
}

// GetLessonByID обрабатывает получение урока по ID
func (h *LessonHandler) GetLessonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	lesson, err := h.lessonService.GetLessonByID(id)
	if err != nil {
		if err.Error() == "урок не найден" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Ошибка получения урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}

// GetLessonWithWords обрабатывает получение урока со словами
func (h *LessonHandler) GetLessonWithWords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	lessonWithWords, err := h.lessonService.GetLessonWithWords(id)
	if err != nil {
		if err.Error() == "урок не найден" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Ошибка получения урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessonWithWords)
}

// StartLesson обрабатывает начало урока
func (h *LessonHandler) StartLesson(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из заголовка (временно используем фиксированный ID)
	userID := 1

	vars := mux.Vars(r)
	idStr := vars["id"]

	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	if err := h.userLessonService.StartLesson(userID, lessonID); err != nil {
		log.Printf("Ошибка начала урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Урок начат",
	})
}

// CompleteLesson обрабатывает завершение урока
func (h *LessonHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	userID := 1

	vars := mux.Vars(r)
	idStr := vars["id"]

	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Получаем баллы из тела запроса
	var input struct {
		Score int `json:"score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.userLessonService.CompleteLesson(userID, lessonID, input.Score); err != nil {
		log.Printf("Ошибка завершения урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Урок завершён",
	})
}

// GetUserLessons обрабатывает получение уроков пользователя
func (h *LessonHandler) GetUserLessons(w http.ResponseWriter, r *http.Request) {
	userID := 1

	userLessons, err := h.userLessonService.GetUserLessons(userID)
	if err != nil {
		log.Printf("Ошибка получения уроков пользователя: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userLessons)
}

// GetProgressStats обрабатывает получение статистики прогресса
func (h *LessonHandler) GetProgressStats(w http.ResponseWriter, r *http.Request) {
	userID := 1

	stats, err := h.userLessonService.GetProgressStats(userID)
	if err != nil {
		log.Printf("Ошибка получения статистики: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetTeacherLessons обрабатывает получение уроков учителя
func (h *LessonHandler) GetTeacherLessons(w http.ResponseWriter, r *http.Request) {
	teacherID := 1 // Временно используем фиксированный ID

	lessons, err := h.lessonService.GetTeacherLessons(teacherID)
	if err != nil {
		log.Printf("Ошибка получения уроков учителя: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons)
}

// CreateLesson обрабатывает создание урока
func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	teacherID := 1

	var input models.LessonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Название урока обязательно", http.StatusBadRequest)
		return
	}

	lesson, err := h.lessonService.CreateLesson(input, teacherID)
	if err != nil {
		log.Printf("Ошибка создания урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lesson)
}

// UpdateLesson обрабатывает обновление урока
func (h *LessonHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	teacherID := 1

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var input models.LessonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.lessonService.UpdateLesson(id, input, teacherID); err != nil {
		log.Printf("Ошибка обновления урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Урок успешно обновлен",
	})
}

// DeleteLesson обрабатывает удаление урока
func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	if err := h.lessonService.DeleteLesson(id); err != nil {
		log.Printf("Ошибка удаления урока: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Урок успешно удален",
	})
}
