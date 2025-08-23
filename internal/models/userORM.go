package models

import (
	"time"

	"gorm.io/gorm"
)

// User представляет модель пользователя в базе данных
type User struct {
	ID        string         `json:"id" gorm:"primary_key"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Tasks     []Task         `json:"tasks" gorm:"foreignkey:UserID;references:ID"` // Связь с задачами
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRequest Используется при обработке входящих запросов
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse Используется для отправки данных клиенту
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"` // Конфиденциальное поле (исключительно для отладки)
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserUpdate представляет структуру для частичного обновления пользователя
type UserUpdate struct {
	Email    *string `json:"email,omitempty"`    // Новый email (опционально)
	Password *string `json:"password,omitempty"` // Новый пароль (опционально)
}
