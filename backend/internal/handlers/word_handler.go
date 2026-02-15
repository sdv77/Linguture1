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

// WordHandler обрабатывает HTTP запросы
type WordHandler struct {
	service *service.WordService
}

// NewWordHandler создает новый обработчик
func NewWordHandler(service *service.WordService) *WordHandler {
	return &WordHandler{service: service}
}

// CreateWord обрабатывает создание слова
func (h *WordHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из заголовка (добавленного в middleware)
	// В реальном приложении используем контекст
	// Здесь для простоты передадим фиктивный ID
	userID := 1 // TODO: Получить реальный userID из токена

	var input models.WordInput

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if input.Word == "" || input.Meaning == "" {
		http.Error(w, "Поля 'word' и 'meaning' обязательны", http.StatusBadRequest)
		return
	}

	// Создаем слово через сервис
	word, err := h.service.CreateWord(userID, input)
	if err != nil {
		log.Printf("Ошибка создания слова: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный объект
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(word)
}

// GetAllWords обрабатывает получение всех слов
func (h *WordHandler) GetAllWords(w http.ResponseWriter, r *http.Request) {
	userID := 1 // TODO: Получить реальный userID из токена

	words, err := h.service.GetAllWords(userID)
	if err != nil {
		log.Printf("Ошибка получения слов: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}

// GetWordByID обрабатывает получение слова по ID
func (h *WordHandler) GetWordByID(w http.ResponseWriter, r *http.Request) {
	userID := 1 // TODO: Получить реальный userID из токена

	// Получаем ID из переменных пути (path variables)
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	word, err := h.service.GetWordByID(id, userID)
	if err != nil {
		if err.Error() == "слово не найдено" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Ошибка получения слова: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(word)
}

// UpdateWord обрабатывает обновление слова
func (h *WordHandler) UpdateWord(w http.ResponseWriter, r *http.Request) {
	userID := 1 // TODO: Получить реальный userID из токена

	// Получаем ID из переменных пути
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var input models.WordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if input.Word == "" || input.Meaning == "" {
		http.Error(w, "Поля 'word' и 'meaning' обязательны", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateWord(id, userID, input); err != nil {
		if err.Error() == "слово не найдено" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Ошибка обновления слова: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Слово успешно обновлено"})
}

// DeleteWord обрабатывает удаление слова
func (h *WordHandler) DeleteWord(w http.ResponseWriter, r *http.Request) {
	userID := 1 // TODO: Получить реальный userID из токена

	// Получаем ID из переменных пути
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteWord(id, userID); err != nil {
		if err.Error() == "слово не найдено" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Ошибка удаления слова: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Слово успешно удалено"})
}
