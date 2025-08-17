package taskService

import (
	"fmt"

	"gorm.io/gorm"
)

// TaskRepository Интерфейс репозитория для работы с задачами CRUD
type TaskRepository interface {
	GetAll() ([]Task, error)
	GetByID(id string) (Task, error)
	Create(task Task) (Task, error)
	Update(task Task) (Task, error)
	Delete(id string) error
}

// Структура, которая реализует все методы TaskRepository
type taskRepository struct { // Место для таски
	db *gorm.DB // Инструмент подключения к БД
}

// NewTaskRepository Конструктор репозитория
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db} // возвращаем из функции: заворачиваем taskRepository в TaskRepository
}

// GetAll Извлекаем все неудаленные таски из БД
func (r *taskRepository) GetAll() ([]Task, error) {
	// 1. Всегда начинаем с инициализированного слайса
	tasks := make([]Task, 0)

	// 2. Выполняем запрос
	result := r.db.Where("deleted_at IS NULL").Find(&tasks)

	// 3. Обрабатываем ошибки
	if result.Error != nil {
		return nil, fmt.Errorf("repo: could not get all tasks: %w", result.Error)
	}

	// 4. Возвращаем результат (даже если он пустой)
	return tasks, nil
}

// GetByID Поиск задачи по ID
func (r *taskRepository) GetByID(id string) (Task, error) {
	var task Task //место, чтобы временно разместить таску из БД

	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

// Create Создание задачи
func (r *taskRepository) Create(task Task) (Task, error) {
	err := r.db.Create(&task).Error
	return task, err
}

// Update Редактирование задачи
func (r *taskRepository) Update(task Task) (Task, error) {
	err := r.db.Save(&task).Error
	return task, err
}

// Delete Удаление (мягкое) задачи
func (r *taskRepository) Delete(id string) error {
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).Delete(&Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
