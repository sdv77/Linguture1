package admin

import (
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	adminHandlers "github.com/sdv77/Linguture1/internal/admin/handlers"
	adminMiddleware "github.com/sdv77/Linguture1/internal/admin/middleware"
	adminRepo "github.com/sdv77/Linguture1/internal/admin/repository"
	adminService "github.com/sdv77/Linguture1/internal/admin/service"
	"github.com/sdv77/Linguture1/pkg/database"
	"github.com/sdv77/Linguture1/pkg/token"
)

//go:embed admin_panel.html
var embedFiles embed.FS

func Main() {
	// Загружаем .env файл
	godotenv.Load("../../.env")

	// Настройки базы данных
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "words_db")

	// Настройки JWT
	jwtSecret := getEnv("JWT_SECRET", "admin-secret-key-change-in-production")
	jwtExpiry := getEnv("JWT_EXPIRY", "24h")

	tokenExpiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		log.Fatalf("Ошибка парсинга времени жизни токена: %v", err)
	}

	// Подключение к базе данных
	db, err := database.ConnectDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создание сервисов
	tokenService := token.NewService(jwtSecret, tokenExpiry)
	adminRepository := adminRepo.NewAdminRepository(db)
	adminSvc := adminService.NewAdminService(adminRepository, adminRepository, tokenService)

	// Создание хендлеров
	authHandler := adminHandlers.NewAdminAuthHandler(adminSvc)
	backupHandler := adminHandlers.NewBackupHandler(adminSvc, db)
	teacherHandler := adminHandlers.NewTeacherHandler(adminSvc)

	// Middleware
	authMiddleware := adminMiddleware.NewAdminAuthMiddleware(tokenService)

	// Роутер
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	// Публичные маршруты
	router.HandleFunc("/api/admin/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Защищенные маршруты
	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(authMiddleware.Middleware)

	// Бэкапы
	adminRouter.HandleFunc("/backup", backupHandler.CreateBackup).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/restore", backupHandler.RestoreBackup).Methods("POST", "OPTIONS")

	// Управление учителями
	adminRouter.HandleFunc("/teachers", teacherHandler.GetAllTeachers).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/teachers", teacherHandler.CreateTeacher).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", teacherHandler.GetTeacherByID).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", teacherHandler.UpdateTeacher).Methods("PUT", "OPTIONS")
	adminRouter.HandleFunc("/teachers/{id}", teacherHandler.DeleteTeacher).Methods("DELETE", "OPTIONS")

	// Страница админ панели
	router.HandleFunc("/", serveAdminPanel).Methods("GET")

	// Ссылка на Adminer
	router.HandleFunc("/adminer", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8082", http.StatusFound)
	}).Methods("GET")

	port := getEnv("ADMIN_PORT", "8081")
	log.Printf("✅ Админ панель запущена на порту %s", port)
	log.Printf("🌐 Доступ: http://localhost:%s", port)
	log.Printf("👤 Тестовый админ: admin / admin123")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

func serveAdminPanel(w http.ResponseWriter, r *http.Request) {
	html, err := embedFiles.ReadFile("admin_panel.html")
	if err != nil {
		log.Printf("❌ Ошибка чтения HTML: %v", err)
		http.Error(w, "Admin panel HTML not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}
