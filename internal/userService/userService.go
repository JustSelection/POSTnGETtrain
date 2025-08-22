package userService

import (
	"errors"

	"github.com/google/uuid"
)

// Глобальные ошибки сервиса
var (
	ErrEmailExists  = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
)

// UserService Интерфейс сервиса для работы с пользователями
type UserService interface {
	GetAllUsers() ([]User, error)
	CreateUser(email, password string) (*User, error)
	UpdateUser(id string, email, password *string) (*User, error)
	DeleteUser(id string) error
	GetUserByID(id string) (*User, error)
}

// Реализация UserService
type userService struct {
	repo UserRepository // Репозиторий для работы с базой данных
}

// NewUserService Конструктор сервиса
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAllUsers Получение всех пользователей
func (s *userService) GetAllUsers() ([]User, error) {
	return s.repo.GetAll() // Просто делегируем запрос в репозиторий
}

// CreateUser Создание пользователя
func (s *userService) CreateUser(email, password string) (*User, error) {
	user := &User{
		ID:       uuid.New().String(), // Генерируем уникальный ID
		Email:    email,               // Устанавливаем email
		Password: password,            // Устанавливаем пароль
	}
	return s.repo.Create(user) // Передаем создание в репозиторий
}

// UpdateUser Обновление пользователя
func (s *userService) UpdateUser(id string, email, password *string) (*User, error) {
	// Сначала получаем пользователя по ID
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля, если они переданы
	if email != nil {
		user.Email = *email
	}
	if password != nil {
		user.Password = *password
	}

	// Сохраняем изменения через репозиторий
	return s.repo.Update(user)
}

// DeleteUser Удаление пользователя
func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id) // Удаляем через репозиторий
}

// Пока не работает
func (s *userService) GetUserByID(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
