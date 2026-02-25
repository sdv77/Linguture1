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
	// В методе Register, при создании пользователя:
	user := &models.User{
		Email:                    input.Email,
		PasswordHash:             string(passwordHash),
		IsVerified:               false, // по умолчанию не подтверждён
		VerificationToken:        verificationToken,
		VerificationTokenExpires: &expiresAt,

		// 🔹 НОВЫЕ ПОЛЯ: профиль ещё не настроен 🔹
		Nickname:         "",    // пусто, пока пользователь не укажет
		NativeLanguage:   "",    // пусто
		LearningLanguage: "",    // пусто
		IsSetup:          false, // 🔹 ключевое: флаг "нужна настройка"

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		ID:               user.ID,
		Email:            user.Email,
		Nickname:         user.Nickname,         // 🔹 новое
		NativeLanguage:   user.NativeLanguage,   // 🔹 новое
		LearningLanguage: user.LearningLanguage, // 🔹 новое
		IsVerified:       user.IsVerified,
		IsSetup:          user.IsSetup, // 🔹 критически важно для фронтенда!
		CreatedAt:        user.CreatedAt,
	}, nil
}

// Возвращает ошибку, если никнейм уже занят или данные невалидны
// SetupProfile обновляет профиль пользователя после первичной авторизации
func (s *AuthService) SetupProfile(userID int, input models.UserProfileSetupInput) error {
	// 🔹 Простая валидация языков
	validLanguages := map[string]bool{
		"ru": true, "en": true, "es": true, "de": true, "fr": true, "it": true,
		"pt": true, "zh": true, "ja": true, "ko": true, "ar": true, "tr": true,
	}

	if !validLanguages[input.NativeLanguage] {
		return fmt.Errorf("неподдерживаемый родной язык: %s", input.NativeLanguage)
	}
	if !validLanguages[input.LearningLanguage] {
		return fmt.Errorf("неподдерживаемый язык для изучения: %s", input.LearningLanguage)
	}
	if input.NativeLanguage == input.LearningLanguage {
		return fmt.Errorf("родной язык и язык изучения не могут совпадать")
	}

	// 🔹 Валидация никнейма
	if len(input.Nickname) < 3 || len(input.Nickname) > 30 {
		return fmt.Errorf("никнейм должен быть от 3 до 30 символов")
	}
	for _, r := range input.Nickname {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return fmt.Errorf("никнейм может содержать только латинские буквы, цифры и подчёркивание")
		}
	}

	// 🔹 Обновляем в репозитории
	if err := s.userRepo.UpdateProfileSetup(userID, input.Nickname, input.NativeLanguage, input.LearningLanguage); err != nil {
		return err
	}

	return nil
}

// GetUserByEmail получает пользователя по email (для использования в хендлерах)
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}
