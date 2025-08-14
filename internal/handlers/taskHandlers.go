package handlers

import (
	"POSTnGETtrain/internal/taskService" // Импорт сервисного слоя для работы с задачами
	"net/http"                           // Пакет для работы с HTTP

	"github.com/labstack/echo/v4" // Веб-фреймворк Echo
)

// TaskHandler - структура обработчиков задач
type TaskHandler struct {
	service taskService.TaskService // Сервис для бизнес-логики задач
}

// NewTaskHandler - конструктор для TaskHandler
func NewTaskHandler(s taskService.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// GetListTasks - обработчик GET запроса для получения списка всех задач
func (h *TaskHandler) GetListTasks(c echo.Context) error {
	// Получаем все задачи через сервис
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		// В случае ошибки возвращаем 500 статус
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Error"})
	}
	return c.JSON(http.StatusOK, tasks)
}

// GetTask - обработчик GET запроса для получения одной задачи по ID
func (h *TaskHandler) GetTask(c echo.Context) error {
	// Получаем ID задачи из URL
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "task ID is required"})
	}

	// Получаем задачу через сервис
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Error"})
	}
	return c.JSON(http.StatusOK, task)
}

// PostTask - обработчик POST запроса для создания новой задачи
func (h *TaskHandler) PostTask(c echo.Context) error {
	// Структура для парсинга тела запроса
	var newTask struct {
		Name   string `json:"name"`    // Название задачи
		IsDone bool   `json:"is_done"` // Статус выполнения
	}

	// Заворачиваем тело запроса в структуру
	if err := c.Bind(&newTask); err != nil {
		// В случае ошибки парсинга возвращаем 400 статус
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Создаем задачу через сервис
	task, err := h.service.CreateTask(newTask.Name, newTask.IsDone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Error"})
	}

	return c.JSON(http.StatusCreated, task)
}

// PatchTask - обработчик PATCH запроса для частичного обновления задачи
func (h *TaskHandler) PatchTask(c echo.Context) error {
	// Получаем ID задачи из параметров URL
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "task ID is required"})
	}

	// Структура для тела запроса
	var updates struct {
		Name   *string `json:"name,omitempty"`    // Указатель на новое название (опционально)
		IsDone *bool   `json:"is_done,omitempty"` // Указатель на новый статус (опционально)
	}

	// Заворачиваем тело запроса в структуру
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	// Обновляем задачу через сервис
	task, err := h.service.UpdateTask(id, updates.Name, updates.IsDone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Error"})
	}
	return c.JSON(http.StatusOK, task)
}

// DeleteTask - обработчик DELETE запроса для удаления задачи
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	// Получаем ID задачи из параметров URL
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "task ID is required"})
	}

	// Удаляем задачу через сервис
	if err := h.service.DeleteTask(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database Error"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
