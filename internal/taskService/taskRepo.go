package taskService

import (
	"POSTnGETtrain/internal/models"
	"fmt"

	"gorm.io/gorm"
)

// TaskRepository Интерфейс репозитория для работы с задачами CRUD
type TaskRepository interface {
	GetAll() ([]models.Task, error)
	GetByID(id string) (models.Task, error)
	GetByUserID(userID string) ([]models.Task, error)
	Create(task models.Task) (models.Task, error)
	Update(task models.Task) (models.Task, error)
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
func (r *taskRepository) GetAll() ([]models.Task, error) {
	// Всегда начинаем с инициализированного слайса
	tasks := make([]models.Task, 0)

	// Выполняем запрос
	result := r.db.Where("deleted_at IS NULL").Find(&tasks)

	// Обрабатываем ошибки
	if result.Error != nil {
		return nil, fmt.Errorf("repo: could not get all tasks: %w", result.Error)
	}

	// Возвращаем результат (даже если он пустой)
	return tasks, nil
}

// GetByID Поиск задачи по ID
func (r *taskRepository) GetByID(id string) (models.Task, error) {
	var task models.Task // место, чтобы временно разместить таску из БД

	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&task)
	if result.Error != nil {
		return models.Task{}, fmt.Errorf("repo: could not get task by id: %w", result.Error)
	}
	return task, nil
}

// Create Создание задачи
func (r *taskRepository) Create(task models.Task) (models.Task, error) {
	err := r.db.Create(&task).Error
	return task, err
}

// Update Редактирование задачи
func (r *taskRepository) Update(task models.Task) (models.Task, error) {
	err := r.db.Save(&task).Error
	return task, err
}

// Delete Удаление (мягкое) задачи
func (r *taskRepository) Delete(id string) error {
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *taskRepository) GetByUserID(userID string) ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&tasks)
	if result.Error != nil {
		return nil, fmt.Errorf("repo: could not get tasks for user %s: %w", userID, result.Error)
	}
	return tasks, nil
}
