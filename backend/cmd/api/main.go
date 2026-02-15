package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sdv77/Linguture1/internal/handlers"
	"github.com/sdv77/Linguture1/internal/middleware"
	"github.com/sdv77/Linguture1/internal/repository"
	"github.com/sdv77/Linguture1/internal/service"
	"github.com/sdv77/Linguture1/pkg/database"
	"github.com/sdv77/Linguture1/pkg/email"
	"github.com/sdv77/Linguture1/pkg/token"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	godotenv.Load()

	// Получаем настройки из переменных окружения
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "words_db")

	// Настройки JWT
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	jwtExpiry := getEnv("JWT_EXPIRY", "24h")

	// Настройки почты
	smtpHost := getEnv("SMTP_HOST", "smtp.mail.ru")
	smtpPortStr := getEnv("SMTP_PORT", "465")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	smtpFrom := getEnv("SMTP_FROM", smtpUser)

	// Парсим порт
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Ошибка парсинга порта SMTP: %v", err)
	}

	// Парсим время жизни токена
	tokenExpiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		log.Fatalf("Ошибка парсинга времени жизни токена: %w", err)
	}

	// Подключаемся к базе данных
	db, err := database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем сервис токенов
	tokenService := token.NewService(jwtSecret, tokenExpiry)

	// Создаем почтовый сервис (может быть nil если нет настроек)
	var emailService *email.Service
	if smtpUser != "" && smtpPassword != "" {
		emailConfig := email.Config{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUser,
			Password: smtpPassword,
			From:     smtpFrom,
			// Убрали UseSSL - теперь определяется автоматически по порту
		}
		emailService = email.NewService(emailConfig)
		log.Printf("Почтовый сервис настроен: %s:%d", smtpHost, smtpPort)
	} else {
		log.Println("Предупреждение: почтовый сервис не настроен (нет SMTP_USER или SMTP_PASSWORD)")
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	wordRepo := repository.NewWordRepository(db)

	// Создаем сервисы
	authService := service.NewAuthService(userRepo, tokenService, emailService)
	wordService := service.NewWordService(wordRepo)

	// Создаем хендлеры
	authHandler := handlers.NewAuthHandler(authService, emailService)
	wordHandler := handlers.NewWordHandler(wordService)

	// Создаем middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Создаем роутер
	router := mux.NewRouter()

	// Добавляем обработчик CORS
	router.Use(corsMiddleware)

	// Публичные маршруты (не требуют аутентификации)
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/verify", authHandler.VerifyEmail).Methods("GET", "OPTIONS")

	// Защищенные маршруты (требуют аутентификации)
	wordsRouter := router.PathPrefix("/api/words").Subrouter()
	wordsRouter.Use(authMiddleware.Middleware)
	wordsRouter.HandleFunc("", wordHandler.CreateWord).Methods("POST", "OPTIONS")
	wordsRouter.HandleFunc("", wordHandler.GetAllWords).Methods("GET", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.GetWordByID).Methods("GET", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.UpdateWord).Methods("PUT", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.DeleteWord).Methods("DELETE", "OPTIONS")

	// Маршрут для получения текущего пользователя
	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.Use(authMiddleware.Middleware)
	userRouter.HandleFunc("/me", authHandler.GetCurrentUser).Methods("GET", "OPTIONS")

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
