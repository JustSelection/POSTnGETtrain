package taskService

import (
	"POSTnGETtrain/internal/models"
	"fmt"

	"github.com/google/uuid" // Пакет для генерации UUID
)

// TaskService - интерфейс сервиса для работы с задачами
type TaskService interface {
	GetAllTasks() ([]models.Task, error)                                                   // Получить все задачи
	GetTaskByID(id string) (models.Task, error)                                            // Получить задачу по ID
	CreateTask(name string, isDone bool, userID string) (models.Task, error)               // Создать новую задачу
	UpdateTask(id string, name *string, isDone *bool, userID *string) (models.Task, error) // Обновить задачу
	DeleteTask(id string) error                                                            // Удалить задачу
	GetTasksByUserID(userID string) ([]models.Task, error)
}

// Реализация интерфейса TaskService
type taskService struct {
	repo TaskRepository // Репозиторий для работы с хранилищем данных
}

// NewTaskService Конструктор сервиса задач
func NewTaskService(r TaskRepository) TaskService {
	return &taskService{repo: r} // Возвращаем указатель на созданный сервис
}

// GetAllTasks - получение всех задач
func (s *taskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll() // Получаем список задач через репозиторий
}

// GetTaskByID Получение задачи по идентификатору
func (s *taskService) GetTaskByID(id string) (models.Task, error) {
	return s.repo.GetByID(id) // Получаем задачу через репозиторий
}

// CreateTask Создание новой задачи
func (s *taskService) CreateTask(name string, isDone bool, userID string) (models.Task, error) {

	task := models.Task{
		ID:     uuid.NewString(), // Генерируем новый UUID
		Name:   name,             // Устанавливаем название
		IsDone: isDone,           // Устанавливаем статус
		UserID: userID,           // Принадлежность пользователю
	}
	return s.repo.Create(task) // Сохраняем через репозиторий
}

// UpdateTask Обновление существующей задачи
func (s *taskService) UpdateTask(id string, name *string, isDone *bool, userID *string) (models.Task, error) {
	// Получаем текущую задачу из репозитория
	task, err := s.repo.GetByID(id)
	if err != nil {
		return models.Task{}, err // Возвращаем ошибку если задача не найдена
	}

	// Обновляем название если передан новый параметр
	if name != nil {
		task.Name = *name // Забираем название
	}

	// Обновляем статус если передан новый параметр
	if isDone != nil {
		task.IsDone = *isDone // Забираем статус выполнения
	}

	if userID != nil {
		task.UserID = *userID
	}

	// Сохраняем измененную задачу через репозиторий
	return s.repo.Update(task)
}

// DeleteTask Удаление задачи по ИДу
func (s *taskService) DeleteTask(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("service: could not delete task %s: %w", id, err)
	}
	return nil
}

func (s *taskService) GetTasksByUserID(userID string) ([]models.Task, error) {
	// Метод в репозитории!
	return s.repo.GetByUserID(userID)
}
