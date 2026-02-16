package handlers

import (
	"database/sql" // ДОБАВЛЕНО: для работы с *sql.DB
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sdv77/Linguture1/internal/admin/service"
)

type BackupHandler struct {
	adminService *service.AdminService
	db           *sql.DB // ДОБАВЛЕНО: поле для базы данных
}

// ИСПРАВЛЕНО: Теперь принимает 2 аргумента
func NewBackupHandler(adminService *service.AdminService, db *sql.DB) *BackupHandler {
	return &BackupHandler{
		adminService: adminService,
		db:           db,
	}
}

func (h *BackupHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "words_db")

	backupData, err := h.adminService.BackupDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Printf("Ошибка создания бэкапа: %v", err)
		http.Error(w, "Ошибка создания бэкапа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=backup_"+getTimestamp()+".sql")
	w.Header().Set("Content-Type", "application/sql")
	w.Write(backupData)
}

// RestoreBackup восстанавливает базу данных из файла
func (h *BackupHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	// Проверяем тип содержимого
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}

	// Получаем файл (ИСПРАВЛЕНО: убираем неиспользуемую переменную handler)
	file, _, err := r.FormFile("backup") // Заменили handler на _
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// Выполняем SQL (только для разработки!)
	// ВАЖНО: В реальном приложении нужно разбивать на отдельные запросы!
	_, err = h.db.Exec(string(content))
	if err != nil {
		log.Printf("Ошибка восстановления БД: %v", err)
		http.Error(w, "Ошибка восстановления базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "База данных успешно восстановлена! Перезагрузите страницу.",
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getTimestamp() string {
	return time.Now().Format("20060102_150405")
}
