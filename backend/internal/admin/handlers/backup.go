package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	adminService "github.com/sdv77/Linguture1/internal/admin/service"
)

// BackupHandler обработчик для бэкапов
type BackupHandler struct {
	service *adminService.AdminService
	db      *sql.DB
}

// NewBackupHandler создаёт новый BackupHandler
func NewBackupHandler(service *adminService.AdminService, db *sql.DB) *BackupHandler {
	return &BackupHandler{
		service: service,
		db:      db,
	}
}

// CreateBackup создаёт SQL-дамп базы данных и отдаёт его для скачивания
func (h *BackupHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Формируем имя файла с текущей датой и временем
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)

	// 🔹 Команда pg_dump для создания дампа
	// Используем переменные окружения для подключения
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbName := getEnv("DB_NAME", "words_db")
	dbPassword := getEnv("DB_PASSWORD", "postgres")

	cmd := exec.Command("pg_dump",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"--format=plain", // Обычный SQL файл
		"--no-owner",     // Не включать владельца
		"--no-acl",       // Не включать права доступа
		"--clean",        // Добавлять команды DROP TABLE
		"--if-exists",    // Проверять существование перед удалением
	)

	// Передаём пароль через окружение (безопаснее чем в аргументах)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

	// Буфер для вывода команды
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Запускаем команду
	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Backup failed: %s", stderr.String()), http.StatusInternalServerError)
		return
	}

	// Отправляем файл браузеру
	w.Header().Set("Content-Type", "application/sql")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", out.Len()))
	w.Write(out.Bytes())
}

// RestoreBackup восстанавливает базу данных из SQL-файла
func (h *BackupHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Парсим multipart/form-data (лимит 100MB)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, handler, err := r.FormFile("backup")
	if err != nil {
		http.Error(w, "No backup file provided. Use field name 'backup'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("📥 Загружен файл: %s, размер: %d bytes\n", handler.Filename, handler.Size)

	// Создаём временную папку для бэкапов (если нет)
	backupDir := "/app/backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		http.Error(w, "Failed to create backup directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем файл во временную папку
	tempPath := filepath.Join(backupDir, fmt.Sprintf("restore_%s.sql", time.Now().Format("150405")))
	outFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		outFile.Close()
		os.Remove(tempPath) // Удаляем после восстановления
	}()

	// Копируем содержимое загруженного файла
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to write file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("💾 Временный файл сохранён: %s\n", tempPath)

	// 🔹 Команда psql для восстановления
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbName := getEnv("DB_NAME", "words_db")
	dbPassword := getEnv("DB_PASSWORD", "postgres")

	cmd := exec.Command("psql",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-f", tempPath,
		"--single-transaction",      // Всё в одной транзакции
		"--set", "ON_ERROR_STOP=on", // Остановиться при первой ошибке
		"-v", "ON_ERROR_STOP=1", // Альтернативный способ
	)

	// Передаём пароль через окружение
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Запускаем восстановление
	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Restore failed: %s", stderr.String()), http.StatusInternalServerError)
		return
	}

	// ✅ Успех
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Database restored successfully", "status": "ok"}`))
}

// Вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
