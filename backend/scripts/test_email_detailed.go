package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sdv77/Linguture1/pkg/email"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️  Файл .env не найден")
	}

	// Получаем настройки
	smtpHost := getEnv("SMTP_HOST", "smtp.mail.ru")
	smtpPortStr := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	smtpFrom := getEnv("SMTP_FROM", smtpUser)

	// Проверяем обязательные поля
	if smtpUser == "" {
		log.Fatal("❌ Ошибка: не задана переменная SMTP_USER")
	}
	if smtpPassword == "" {
		log.Fatal("❌ Ошибка: не задана переменная SMTP_PASSWORD")
	}

	// Парсим порт
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("❌ Ошибка парсинга порта: %v", err)
	}

	fmt.Println("📧 Настройки почты:")
	fmt.Printf("   SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("   SMTP_PORT: %d\n", smtpPort)
	fmt.Printf("   SMTP_USER: %s\n", smtpUser)
	fmt.Printf("   SMTP_PASSWORD: %s\n", maskPassword(smtpPassword))
	fmt.Printf("   SMTP_FROM: %s\n", smtpFrom)
	fmt.Println()

	// Проверяем доступность сервера
	fmt.Printf("🔍 Проверка доступности %s:%d...\n", smtpHost, smtpPort)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), 10*time.Second)
	if err != nil {
		log.Fatalf("❌ Сервер недоступен: %v\nПроверьте:\n1. Подключение к интернету\n2. Не блокирует ли провайдер порт %d", err, smtpPort)
	}
	conn.Close()
	fmt.Println("✅ Сервер доступен")
	fmt.Println()

	// Создаём почтовый сервис
	cfg := email.Config{
		Host:     smtpHost,
		Port:     smtpPort,
		Username: smtpUser,
		Password: smtpPassword,
		From:     smtpFrom,
	}

	svc := email.NewService(cfg)

	// Тестовый адрес
	testEmail := smtpUser

	fmt.Printf("📤 Отправка тестового письма на %s...\n", testEmail)
	fmt.Println()

	start := time.Now()
	err = svc.SendVerificationEmail(testEmail, "test-token-123456")
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ Ошибка отправки (время: %v): %v", duration, err)

		// Подсказки по решению
		if smtpPort == 465 {
			fmt.Println("\n💡 Попробуйте:")
			fmt.Println("1. Используйте порт 587 вместо 465")
			fmt.Println("2. Убедитесь, что используете пароль для внешних приложений из настроек mail.ru")
			fmt.Println("3. Проверьте, включена ли опция 'Пароли для внешних приложений' в mail.ru")
		}

		log.Fatal("Отправка не удалась")
	}

	fmt.Printf("✅ Письмо отправлено успешно! (время: %v)\n", duration)
	fmt.Println("   Проверьте почтовый ящик (и папку Спам)")
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// maskPassword маскирует пароль для вывода
func maskPassword(p string) string {
	if len(p) <= 4 {
		return "***"
	}
	return p[:2] + "***" + p[len(p)-2:]
}
