package handlers

import (
	"POSTnGETtrain/internal/taskService"
	"POSTnGETtrain/internal/web/tasks"
	"context"
	"fmt"
)

// Handler - заготовка для конструктора
type Handler struct {
	service taskService.TaskService // Сервис для бизнес-логики работы с задачами
}

// NewHandler - сам конструктор
func NewHandler(s taskService.TaskService) *Handler {
	return &Handler{service: s}
}

// GetTasks - ctx и request не используются, но требуются для запроса
func (h *Handler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (
	tasks.GetTasksResponseObject, error) {

	// Получаем все задачи из сервисного слоя
	dbTasks, err := h.service.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("handler: could not get all tasks: %w", err)
	}

	// Создаем пустой список для ответа
	response := make(tasks.GetTasks200JSONResponse, 0)

	// Преобразуем задачи из формата сервиса в формат API
	for _, t := range dbTasks {
		response = append(response, tasks.Task{
			ID:     t.ID,     // Идентификатор задачи
			Name:   t.Name,   // Название задачи
			IsDone: t.IsDone, // Статус выполнения
			UserID: t.UserID, // Какому пользователю принадлежит
		})
	}
	return response, nil // Возвращаем список задач
}

func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (
	tasks.PostTasksResponseObject, error) {
	if request.Body.UserID == "" {
		return nil, fmt.Errorf("handler: user_id is required")
	}
	// Устанавливаем статус по умолчанию
	isDone := false
	if request.Body.IsDone != nil {
		isDone = *request.Body.IsDone // Если статус указан, используем его
	}

	// Создаем задачу с запросом в сервис
	created, err := h.service.CreateTask(request.Body.Name, isDone, request.Body.UserID)
	if err != nil {
		return nil, fmt.Errorf("handler: could not create task: %w", err) // Обрабатываем ошибку создания
	}

	// Возвращаем созданную задачу в формате API
	return tasks.PostTasks201JSONResponse{
		ID:     created.ID,     // ID созданной задачи
		Name:   created.Name,   // Название задачи
		IsDone: created.IsDone, // Статус выполнения
		UserID: created.UserID,
	}, nil
}

// GetUsersIdTasks - получить все задачи юзера
func (h *Handler) GetUsersIdTasks(_ context.Context, request tasks.GetUsersIdTasksRequestObject) (tasks.GetUsersIdTasksResponseObject, error) {
	tasksList, err := h.service.GetTasksByUserID(request.Id)
	if err != nil {
		return nil, fmt.Errorf("handler: could not get tasks for user%s: %w", request.Id, err)
	}

	response := make(tasks.GetUsersIdTasks200JSONResponse, 0)
	for _, t := range tasksList {
		response = append(response, tasks.Task{
			ID:     t.ID,
			Name:   t.Name,
			IsDone: t.IsDone,
			UserID: t.UserID,
		})
	}
	return response, nil
}

func (h *Handler) GetTasksId(_ context.Context, request tasks.GetTasksIdRequestObject) (
	tasks.GetTasksIdResponseObject, error) {
	// Получаем задачу из сервиса по ID
	task, err := h.service.GetTaskByID(request.Id)
	if err != nil {
		return nil, fmt.Errorf("handler: could not get task by ID %s: %w", request.Id, err) // Обрабатываем ошибку поиска
	}

	// Возвращаем найденную задачу
	return tasks.GetTasksId200JSONResponse{
		ID:     task.ID,     // ID задачи
		Name:   task.Name,   // Название задачи
		IsDone: task.IsDone, // Статус выполнения
		UserID: task.UserID,
	}, nil
}

func (h *Handler) PatchTasksId(_ context.Context, request tasks.PatchTasksIdRequestObject) (
	tasks.PatchTasksIdResponseObject, error) {
	// Переменные для обновляемых полей
	var name *string
	var isDone *bool
	var userID *string

	// Если в запросе указано новое название, сохраняем его
	if request.Body.Name != nil {
		name = request.Body.Name
	}
	// Если в запросе указан новый статус, сохраняем его
	if request.Body.IsDone != nil {
		isDone = request.Body.IsDone
	}

	if request.Body.UserID != nil {
		userID = request.Body.UserID
	}

	// Обновляем задачу через сервис
	updated, err := h.service.UpdateTask(request.Id, name, isDone, userID)
	if err != nil {
		return nil, fmt.Errorf("handler: could not update task %s: %w", request.Id, err) // Обрабатываем ошибку обновления
	}

	// Возвращаем обновленную задачу
	return tasks.PatchTasksId200JSONResponse{
		ID:     updated.ID,     // ID задачи
		Name:   updated.Name,   // Новое название
		IsDone: updated.IsDone, // Новый статус
		UserID: updated.UserID,
	}, nil
}

func (h *Handler) DeleteTasksId(_ context.Context, request tasks.DeleteTasksIdRequestObject) (
	tasks.DeleteTasksIdResponseObject, error) {
	// Удаляем задачу через сервис
	if err := h.service.DeleteTask(request.Id); err != nil {
		return nil, fmt.Errorf("handler: could not delete task %s: %w", request.Id, err)
	}
	return tasks.DeleteTasksId204Response{}, nil
}
