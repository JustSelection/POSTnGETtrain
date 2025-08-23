package models

import "gorm.io/gorm"

type Task struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	IsDone    bool           `json:"is_done"`
	UserID    string         `json:"user_id" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // "-", чтобы это техническое поле не отображалось в JSON
}

// Реализация TaskReference для User

func (t Task) GetID() string     { return t.ID }
func (t Task) GetName() string   { return t.Name }
func (t Task) GetIsDone() bool   { return t.IsDone }
func (t Task) GetUserID() string { return t.UserID }

type TaskRequest struct {
	Name   string `json:"name"`
	IsDone bool   `json:"is_done"`
	UserID string `json:"user_id"`
}
