-- Создаем базу данных (выполняется от пользователя postgres)
CREATE DATABASE words_db;

-- Подключаемся к базе
\c words_db

-- Создаем таблицу слов (если ещё не существует)
CREATE TABLE IF NOT EXISTS words (
    id SERIAL PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    meaning TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_words_word ON words(word);

-- Тестовые данные (если таблица пустая)
INSERT INTO words (word, meaning) VALUES
('hello', 'Приветствие на английском языке'),
('world', 'Планета Земля или общество людей'),
('golang', 'Язык программирования от Google'),
('vue', 'JavaScript фреймворк для создания интерфейсов'),
('postgres', 'Реляционная система управления базами данных')
ON CONFLICT DO NOTHING;