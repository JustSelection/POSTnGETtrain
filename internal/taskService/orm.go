package taskService

import "gorm.io/gorm"

type Task struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // "-", чтобы это техническое поле не отображалось в JSON
}

type TaskRequest struct {
	Name   string `json:"name"`
	IsDone bool   `json:"is_done"`
}
