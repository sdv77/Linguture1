package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sdv77/Linguture1/internal/handlers"
	"github.com/sdv77/Linguture1/internal/repository"
	"github.com/sdv77/Linguture1/internal/service"
	"github.com/sdv77/Linguture1/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	// Получаем настройки из переменных окружения
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "moderator")
	dbPassword := getEnv("DB_PASSWORD", "123")
	dbName := getEnv("DB_NAME", "words_db")

	// Подключаемся к базе данных
	db, err := database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем репозиторий
	wordRepo := repository.NewWordRepository(db)

	// Создаем сервис
	wordService := service.NewWordService(wordRepo)

	// Создаем хендлер
	wordHandler := handlers.NewWordHandler(wordService)

	// Создаем роутер
	router := mux.NewRouter()

	// Добавляем обработчик CORS
	router.Use(corsMiddleware)

	// Определяем маршруты
	// Создание слова
	router.HandleFunc("/api/words", wordHandler.CreateWord).Methods("POST", "OPTIONS")
	// Получение всех слов
	router.HandleFunc("/api/words", wordHandler.GetAllWords).Methods("GET", "OPTIONS")
	// Получение слова по ID
	router.HandleFunc("/api/words/{id}", wordHandler.GetWordByID).Methods("GET", "OPTIONS")
	// Обновление слова
	router.HandleFunc("/api/words/{id}", wordHandler.UpdateWord).Methods("PUT", "OPTIONS")
	// Удаление слова
	router.HandleFunc("/api/words/{id}", wordHandler.DeleteWord).Methods("DELETE", "OPTIONS")

	// Запускаем сервер
	port := getEnv("PORT", "8080")
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// corsMiddleware добавляет заголовки CORS для разрешения запросов с фронтенда
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника (для разработки)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешаем методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Разрешаем заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Разрешаем credentials (если понадобится)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка preflight запросов
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
