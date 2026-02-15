-- Создаем базу данных (если не существует)
SELECT 'CREATE DATABASE words_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'words_db')\gexec

-- Подключаемся к базе
\c words_db

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    verification_token_expires TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу учителей
CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу уроков
CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL DEFAULT 1,
    lesson_type VARCHAR(50) DEFAULT 'vocabulary',
    is_active BOOLEAN DEFAULT TRUE,
    order_num INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем столбец teacher_id если его нет
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'lessons' AND column_name = 'teacher_id') THEN
        ALTER TABLE lessons ADD COLUMN teacher_id INTEGER REFERENCES teachers(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Создаем таблицу слов в уроках (общие слова для всех)
CREATE TABLE IF NOT EXISTS lesson_vocabulary (
    id SERIAL PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    meaning TEXT NOT NULL,
    transcription VARCHAR(255),
    example TEXT,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу прогресса пользователей по урокам
CREATE TABLE IF NOT EXISTS user_lessons (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'not_started',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    score INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, lesson_id)
);

-- Создаем таблицу слов пользователей (личные слова)
CREATE TABLE IF NOT EXISTS words (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    meaning TEXT NOT NULL,
    source_lesson_id INTEGER REFERENCES lessons(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);
CREATE INDEX IF NOT EXISTS idx_lessons_teacher ON lessons(teacher_id);
CREATE INDEX IF NOT EXISTS idx_lessons_level ON lessons(level);
CREATE INDEX IF NOT EXISTS idx_lessons_order ON lessons(order_num);
CREATE INDEX IF NOT EXISTS idx_user_lessons_user ON user_lessons(user_id);
CREATE INDEX IF NOT EXISTS idx_user_lessons_lesson ON user_lessons(lesson_id);
CREATE INDEX IF NOT EXISTS idx_words_user_id ON words(user_id);
CREATE INDEX IF NOT EXISTS idx_words_word ON words(word);
CREATE INDEX IF NOT EXISTS idx_lesson_vocabulary_lesson ON lesson_vocabulary(lesson_id);

-- Тестовый учитель (пароль: teacher123)
INSERT INTO teachers (email, password, full_name, is_active) VALUES
('teacher@wordsapp.com', 'teacher123', 'Иван Петров', TRUE)
ON CONFLICT (email) DO NOTHING;

-- Обновляем существующие уроки, устанавливая teacher_id = 1 (первый учитель)
UPDATE lessons SET teacher_id = 1 WHERE teacher_id IS NULL;

-- Тестовые уроки (только если таблица пустая)
INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Основные приветствия', 'Изучите основные приветствия на английском языке', 1, 'vocabulary', 1, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Основные приветствия');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Семья и родственники', 'Слова о семье и родственниках', 1, 'vocabulary', 2, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Семья и родственники');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Еда и напитки', 'Основные слова о еде и напитках', 2, 'vocabulary', 3, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Еда и напитки');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Цвета', 'Изучите названия цветов', 1, 'vocabulary', 4, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Цвета');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Числа 1-10', 'Основные числа', 1, 'vocabulary', 5, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Числа 1-10');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Дни недели', 'Названия дней недели', 2, 'vocabulary', 6, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Дни недели');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Погода', 'Слова о погоде', 3, 'vocabulary', 7, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Погода');

INSERT INTO lessons (title, description, level, lesson_type, order_num, is_active, teacher_id) 
SELECT 'Профессии', 'Различные профессии', 3, 'vocabulary', 8, TRUE, 1
WHERE NOT EXISTS (SELECT 1 FROM lessons WHERE title = 'Профессии');

-- Тестовые слова для первого урока
INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'hello', 'привет', '[həˈləʊ]', 'Hello, how are you?', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'hello' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'hi', 'привет (неформально)', '[haɪ]', 'Hi there!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'hi' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'goodbye', 'до свидания', '[ˌɡʊdˈbaɪ]', 'Goodbye, see you later!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'goodbye' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'bye', 'пока (неформально)', '[baɪ]', 'Bye bye!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'bye' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'good morning', 'доброе утро', '[ˌɡʊd ˈmɔːnɪŋ]', 'Good morning!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'good morning' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'good afternoon', 'добрый день', '[ˌɡʊd ˌɑːftəˈnuːn]', 'Good afternoon!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'good afternoon' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'good evening', 'добрый вечер', '[ˌɡʊd ˈiːvnɪŋ]', 'Good evening!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'good evening' AND lesson_id = 1);

INSERT INTO lesson_vocabulary (word, meaning, transcription, example, lesson_id) 
SELECT 'good night', 'спокойной ночи', '[ˌɡʊd ˈnaɪt]', 'Good night!', 1
WHERE NOT EXISTS (SELECT 1 FROM lesson_vocabulary WHERE word = 'good night' AND lesson_id = 1);