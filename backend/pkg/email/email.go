package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// Config конфигурация почтового сервера
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Service отправляет письма
type Service struct {
	config Config
}

// NewService создает новый почтовый сервис
func NewService(cfg Config) *Service {
	return &Service{config: cfg}
}

// SendVerificationEmail отправляет письмо с подтверждением
func (s *Service) SendVerificationEmail(to, token string) error {
	// Формируем тело письма
	subject := "Подтверждение регистрации"
	verifyURL := fmt.Sprintf("http://localhost:8080/api/auth/verify?token=%s", token)

	// HTML тело письма
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Подтверждение регистрации</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
        <div style="background: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0;">
            <h1>Добро пожаловать!</h1>
        </div>
        <div style="background: white; padding: 30px; border: 1px solid #ddd; border-top: none; border-radius: 0 0 8px 8px;">
            <p>Спасибо за регистрацию в нашем приложении.</p>
            <p>Пожалуйста, подтвердите вашу почту, перейдя по ссылке ниже:</p>
            <p style="text-align: center; margin: 30px 0;">
                <a href="%s" style="background-color: #4CAF50; color: white; padding: 14px 32px; text-decoration: none; border-radius: 6px; font-size: 18px; font-weight: bold; display: inline-block;">Подтвердить почту</a>
            </p>
            <p>Или скопируйте ссылку:</p>
            <p style="background-color: #f8f9fa; padding: 12px; border-radius: 6px; word-break: break-all; font-family: monospace; font-size: 14px;">
                %s
            </p>
            <p style="margin-top: 25px; color: #666; font-size: 14px;">
                ⏱️ Ссылка действительна 24 часа.<br>
                Если вы не регистрировались, просто проигнорируйте это письмо.
            </p>
        </div>
    </body>
</html>
`, verifyURL, verifyURL)

	// Формируем полное сообщение
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	var conn net.Conn
	var err error

	// Для порта 465 используем SSL с самого начала
	if s.config.Port == 465 {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), &tls.Config{
			ServerName: s.config.Host,
		})
	} else {
		// Для порта 587 используем обычное подключение + STARTTLS
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	}

	if err != nil {
		return fmt.Errorf("ошибка подключения к почтовому серверу: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("ошибка создания клиента SMTP: %w", err)
	}
	defer c.Close()

	// Для порта 587 включаем STARTTLS
	if s.config.Port == 587 {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err = c.StartTLS(&tls.Config{ServerName: s.config.Host}); err != nil {
				return fmt.Errorf("ошибка STARTTLS: %w", err)
			}
		}
	}

	// Аутентификация
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("ошибка аутентификации SMTP: %w. Проверьте логин/пароль и настройки 'Пароли для внешних приложений' в mail.ru", err)
	}

	// Устанавливаем отправителя и получателя
	if err = c.Mail(s.config.From); err != nil {
		return fmt.Errorf("ошибка установки отправителя: %w", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("ошибка установки получателя: %w", err)
	}

	// Отправляем данные
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("ошибка начала передачи данных: %w", err)
	}

	// Формируем заголовки и тело письма
	_, err = fmt.Fprintf(w, "From: %s\r\n", s.config.From)
	if err != nil {
		return fmt.Errorf("ошибка записи заголовка From: %w", err)
	}

	_, err = fmt.Fprintf(w, "To: %s\r\n", to)
	if err != nil {
		return fmt.Errorf("ошибка записи заголовка To: %w", err)
	}

	_, err = fmt.Fprintf(w, "Subject: %s\r\n", subject)
	if err != nil {
		return fmt.Errorf("ошибка записи заголовка Subject: %w", err)
	}

	_, err = fmt.Fprintf(w, "MIME-Version: 1.0\r\n")
	if err != nil {
		return fmt.Errorf("ошибка записи MIME-Version: %w", err)
	}

	_, err = fmt.Fprintf(w, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
	if err != nil {
		return fmt.Errorf("ошибка записи Content-Type: %w", err)
	}

	_, err = fmt.Fprintf(w, "\r\n%s", htmlBody)
	if err != nil {
		return fmt.Errorf("ошибка записи тела письма: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия потока данных: %w", err)
	}

	if err = c.Quit(); err != nil {
		return fmt.Errorf("ошибка завершения сессии SMTP: %w", err)
	}

	fmt.Printf("✅ Письмо отправлено на %s\n", to)
	return nil
}
