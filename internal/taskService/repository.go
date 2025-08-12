package taskService

import (
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
type taskRepository struct { // место для таски
	db *gorm.DB // инструмент подключения к БД
}

// Конструктор репозитория
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db} // возвращаем из функции: заворачиваем taskRepository в TaskRepository
}

// Извлекаем все неудаленные таски из БД
func (r *taskRepository) GetAll() ([]Task, error) {
	var tasks []Task                                           // сюда будем класть найденные в БД таски
	err := r.db.Where("deleted_at IS NULL").Find(&tasks).Error // записываем в &tasks с условием, что не удалено
	return tasks, err
}

// Поиск задачи по ID
func (r *taskRepository) GetByID(id string) (Task, error) {
	var task Task //место, чтобы временно разместить таску из БД

	err := r.db.Where("id = ? AND deleted_at IS NULL", id).Find(&task).Error
	return task, err
}

// Создание задачи
func (r *taskRepository) Create(task Task) (Task, error) {
	err := r.db.Create(&task).Error
	return task, err
}

// Редактирование задачи
func (r *taskRepository) Update(task Task) (Task, error) {
	err := r.db.Save(&task).Error
	return task, err
}

// Удаление (мягкое) задачи
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
