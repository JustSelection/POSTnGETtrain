package userService

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository Содержит все необходимые методы для CRUD операций
type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id string) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(id string) error
	EmailExists(email string) (bool, error)
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
	err := r.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Create создает нового пользователя в базе данных с уникальным email
func (r *userRepository) Create(user *User) (*User, error) {
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
func (r *userRepository) GetByID(id string) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// GetAll возвращает список всех пользователей в системе
func (r *userRepository) GetAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

// Update обновляет данные пользователя в базе данных
func (r *userRepository) Update(user *User) (*User, error) {
	err := r.db.Save(user).Error
	return user, err
}

// Delete удаляет пользователя по его идентификатору
func (r *userRepository) Delete(id string) error {
	result := r.db.Delete(&User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}
