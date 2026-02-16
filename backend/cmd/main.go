// backend/cmd/main.go
package main

import (
	"log"
	"os"

	// Импортируем пакеты admin и api (пути могут отличаться в твоём проекте)

	admin "github.com/sdv77/Linguture1/cmd/admin"

	api "github.com/sdv77/Linguture1/cmd/api"
)

func main() {
	// Проверяем аргумент командной строки
	// Пример: ./app admin или ./app api
	if len(os.Args) < 2 {
		log.Fatal("Укажите режим: 'admin' или 'api'")
	}

	mode := os.Args[1]

	switch mode {
	case "admin":
		log.Println("Запуск Admin Panel...")
		admin.Main() // Вызываем функцию Main из пакета admin
	case "api":
		log.Println("Запуск Backend API...")
		api.Main() // Вызываем функцию Main из пакета api
	default:
		log.Fatalf("Неизвестный режим: %s (доступны: admin, api)", mode)
	}
}
