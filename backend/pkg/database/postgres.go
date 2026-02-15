package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// ConnectDB подключается к базе данных и возвращает соединение
func ConnectDB(host, port, user, password, dbname string) (*sql.DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Открываем соединение с базой
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения: %w", err)
	}

	log.Println("Успешное подключение к базе данных PostgreSQL")

	// Проверяем, существует ли таблица
	if err := checkAndCreateTable(db); err != nil {
		log.Printf("Предупреждение: %v", err)
	}

	return db, nil
}

// checkAndCreateTable проверяет наличие таблицы и создает её при необходимости
func checkAndCreateTable(db *sql.DB) error {
	// Проверяем, существует ли таблица
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'words'
        );
    `).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка проверки таблицы: %w", err)
	}

	if !exists {
		log.Println("Таблица 'words' не найдена, создаём...")

		// Создаем таблицу
		_, err := db.Exec(`
            CREATE TABLE words (
                id SERIAL PRIMARY KEY,
                word VARCHAR(255) NOT NULL,
                meaning TEXT NOT NULL,
                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
            );
        `)

		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}

		log.Println("Таблица 'words' успешно создана")
	} else {
		log.Println("Таблица 'words' уже существует")
	}

	return nil
}
