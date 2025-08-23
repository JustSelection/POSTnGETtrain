package userService

import (
	"POSTnGETtrain/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository Содержит все необходимые методы для CRUD операций
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) error
	EmailExists(email string) (bool, error)
	GetTasksForUser(userID string) ([]models.Task, error)
	GetUserWithTasks(userID string) (*models.User, error)
}

// userRepository - реализация UserRepository с использованием GORM
type userRepository struct {
	db *gorm.DB // Заготовка подключения к базе данных
}

// NewUserRepository создает новый экземпляр userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// EmailExists проверяет, существует ли пользователь с указанным email
func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Create создает нового пользователя в базе данных с уникальным email
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	// Проверяем не занят ли email
	exists, err := r.EmailExists(user.Email)
	if err != nil {
		return nil, fmt.Errorf("email check failed: %w", err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Создаем запись в базе данных
	err = r.db.Create(user).Error
	return user, err
}

// GetByID находит пользователя по ID
func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// GetAll возвращает список всех пользователей в системе
func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update обновляет данные пользователя в базе данных
func (r *userRepository) Update(user *models.User) (*models.User, error) {
	err := r.db.Save(user).Error
	return user, err
}

// Delete удаляет пользователя по его идентификатору
func (r *userRepository) Delete(id string) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

func (r *userRepository) GetTasksForUser(userID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("repo: could not get tasks for user %s: %w", userID, err)
	}
	return tasks, nil
}

func (r *userRepository) GetUserWithTasks(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Tasks", "deleted_at IS NULL").Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: could not get user with tasks %s: %w", userID, err)
	}
	return &user, nil
}
