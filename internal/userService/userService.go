package userService

import (
	"POSTnGETtrain/internal/models"
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
	GetAllUsers() ([]models.User, error)
	CreateUser(email, password string) (*models.User, error)
	UpdateUser(id string, email, password *string) (*models.User, error)
	DeleteUser(id string) error
	GetUserByID(id string) (*models.User, error)
	GetTasksForUser(userID string) ([]models.Task, error)
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
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll() // Просто делегируем запрос в репозиторий
}

// CreateUser Создание пользователя
func (s *userService) CreateUser(email, password string) (*models.User, error) {
	user := &models.User{
		ID:       uuid.New().String(), // Генерируем уникальный ID
		Email:    email,               // Устанавливаем email
		Password: password,            // Устанавливаем пароль
	}
	return s.repo.Create(user) // Передаем создание в репозиторий
}

// UpdateUser Обновление пользователя
func (s *userService) UpdateUser(id string, email, password *string) (*models.User, error) {
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

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetTasksForUser(userID string) ([]models.Task, error) {
	_, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	// Получаем задачи через taskService
	return s.repo.GetTasksForUser(userID)
}
