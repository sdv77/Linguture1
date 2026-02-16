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
	godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "words_db")

	jwtSecret := getEnv("JWT_SECRET", "super-secret-jwt-key-2026")
	jwtExpiry := getEnv("JWT_EXPIRY", "24h")

	smtpHost := getEnv("SMTP_HOST", "smtp.mail.ru")
	smtpPortStr := getEnv("SMTP_PORT", "465")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	smtpFrom := getEnv("SMTP_FROM", smtpUser)

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Ошибка парсинга порта SMTP: %v", err)
	}

	tokenExpiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		log.Fatalf("Ошибка парсинга времени жизни токена: %w", err)
	}

	db, err := database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	tokenService := token.NewService(jwtSecret, tokenExpiry)

	var emailService *email.Service
	if smtpUser != "" && smtpPassword != "" {
		emailConfig := email.Config{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUser,
			Password: smtpPassword,
			From:     smtpFrom,
		}
		emailService = email.NewService(emailConfig)
		log.Printf("Почтовый сервис настроен: %s:%d", smtpHost, smtpPort)
	} else {
		log.Println("Предупреждение: почтовый сервис не настроен")
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	wordRepo := repository.NewWordRepository(db)
	lessonRepo := repository.NewLessonRepository(db)
	userLessonRepo := repository.NewUserLessonRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)

	authService := service.NewAuthService(userRepo, tokenService, emailService)
	wordService := service.NewWordService(wordRepo)
	lessonService := service.NewLessonService(lessonRepo)
	userLessonService := service.NewUserLessonService(userLessonRepo, lessonRepo, wordRepo)
	teacherAuthService := service.NewTeacherAuthService(teacherRepo, tokenService)

	authHandler := handlers.NewAuthHandler(authService, emailService)
	wordHandler := handlers.NewWordHandler(wordService)
	lessonHandler := handlers.NewLessonHandler(lessonService, userLessonService)
	teacherAuthHandler := handlers.NewTeacherAuthHandler(teacherAuthService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	router := mux.NewRouter()

	router.Use(corsMiddleware)

	// Публичные маршруты аутентификации
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/verify", authHandler.VerifyEmail).Methods("GET", "OPTIONS")

	// Защищенные маршруты пользователей
	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.Use(authMiddleware.Middleware)
	userRouter.HandleFunc("/me", authHandler.GetCurrentUser).Methods("GET", "OPTIONS")

	wordsRouter := router.PathPrefix("/api/words").Subrouter()
	wordsRouter.Use(authMiddleware.Middleware)
	wordsRouter.HandleFunc("", wordHandler.CreateWord).Methods("POST", "OPTIONS")
	wordsRouter.HandleFunc("", wordHandler.GetAllWords).Methods("GET", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.GetWordByID).Methods("GET", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.UpdateWord).Methods("PUT", "OPTIONS")
	wordsRouter.HandleFunc("/{id}", wordHandler.DeleteWord).Methods("DELETE", "OPTIONS")

	// Маршруты уроков (публичные и защищенные)
	lessonsPublicRouter := router.PathPrefix("/api/lessons").Subrouter()
	lessonsPublicRouter.HandleFunc("", lessonHandler.GetAllLessons).Methods("GET", "OPTIONS")
	lessonsPublicRouter.HandleFunc("/{id}", lessonHandler.GetLessonByID).Methods("GET", "OPTIONS")
	lessonsPublicRouter.HandleFunc("/{id}/words", lessonHandler.GetLessonWithWords).Methods("GET", "OPTIONS")

	lessonsProtectedRouter := router.PathPrefix("/api/lessons").Subrouter()
	lessonsProtectedRouter.Use(authMiddleware.Middleware)
	lessonsProtectedRouter.HandleFunc("/{id}/start", lessonHandler.StartLesson).Methods("POST", "OPTIONS")
	lessonsProtectedRouter.HandleFunc("/{id}/complete", lessonHandler.CompleteLesson).Methods("POST", "OPTIONS")
	lessonsProtectedRouter.HandleFunc("/my", lessonHandler.GetUserLessons).Methods("GET", "OPTIONS")
	lessonsProtectedRouter.HandleFunc("/stats", lessonHandler.GetProgressStats).Methods("GET", "OPTIONS")

	// Маршруты учителей (аутентификация)
	teacherAuthRouter := router.PathPrefix("/api/teacher/auth").Subrouter()
	teacherAuthRouter.HandleFunc("/login", teacherAuthHandler.Login).Methods("POST", "OPTIONS")
	teacherAuthRouter.HandleFunc("/me", teacherAuthHandler.GetTeacher).Methods("GET", "OPTIONS")

	// Маршруты учителей (уроки)
	teacherLessonsRouter := router.PathPrefix("/api/teacher/lessons").Subrouter()
	teacherLessonsRouter.HandleFunc("", lessonHandler.GetTeacherLessons).Methods("GET", "OPTIONS")
	teacherLessonsRouter.HandleFunc("", lessonHandler.CreateLesson).Methods("POST", "OPTIONS")
	teacherLessonsRouter.HandleFunc("/{id}", lessonHandler.UpdateLesson).Methods("PUT", "OPTIONS")
	teacherLessonsRouter.HandleFunc("/{id}", lessonHandler.DeleteLesson).Methods("DELETE", "OPTIONS")
	teacherLessonsRouter.HandleFunc("/{id}/words", lessonHandler.AddWordToLesson).Methods("POST", "OPTIONS")
	teacherLessonsRouter.HandleFunc("/words/{wordId}", lessonHandler.DeleteWordFromLesson).Methods("DELETE", "OPTIONS")

	port := getEnv("PORT", "8080")
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// corsMiddleware добавляет заголовки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
