package service

import (
	"fmt"
	"time"

	"github.com/sdv77/Linguture1/internal/models"
	"github.com/sdv77/Linguture1/internal/repository"
	"github.com/sdv77/Linguture1/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// AuthService содержит логику аутентификации
type AuthService struct {
	userRepo     *repository.UserRepository
	tokenSvc     *token.Service
	emailService interface { // Добавим интерфейс для отправки почты
		SendVerificationEmail(to, token string) error
	}
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(userRepo *repository.UserRepository, tokenSvc *token.Service, emailService interface {
	SendVerificationEmail(to, token string) error
}) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenSvc:     tokenSvc,
		emailService: emailService,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(input models.UserRegisterInput) error {
	// Проверяем, существует ли пользователь с таким email
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return fmt.Errorf("пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хэширования пароля: %w", err)
	}

	// Генерируем токен подтверждения
	verificationToken, err := s.tokenSvc.GenerateVerificationToken()
	if err != nil {
		return fmt.Errorf("ошибка генерации токена подтверждения: %w", err)
	}

	// Устанавливаем время истечения токена (24 часа)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем пользователя
	user := &models.User{
		Email:                    input.Email,
		PasswordHash:             string(passwordHash),
		VerificationToken:        verificationToken,
		VerificationTokenExpires: &expiresAt,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Отправляем письмо с подтверждением
	// Используем recover для игнорирования ошибок отправки почты
	// (чтобы регистрация прошла даже если почта не настроена)
	if s.emailService != nil {
		go func() {
			if err := s.emailService.SendVerificationEmail(input.Email, verificationToken); err != nil {
				fmt.Printf("Предупреждение: не удалось отправить письмо подтверждения: %v\n", err)
			} else {
				fmt.Printf("Письмо подтверждения отправлено на %s\n", input.Email)
			}
		}()
	}

	return nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(input models.UserLoginInput) (*models.TokenResponse, error) {
	// Находим пользователя по email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Проверяем подтверждён ли email
	if !user.IsVerified {
		return nil, fmt.Errorf("пожалуйста, подтвердите ваш email")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Генерируем JWT токен
	tokenStr, err := s.tokenSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return &models.TokenResponse{
		Token:     tokenStr,
		ExpiresIn: int(s.tokenSvc.GetTokenExpiry().Seconds()),
	}, nil
}

// VerifyEmail подтверждает email пользователя
func (s *AuthService) VerifyEmail(token string) error {
	return s.userRepo.VerifyEmail(token)
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}, nil
}
