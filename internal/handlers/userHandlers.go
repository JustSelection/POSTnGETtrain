package handlers

import (
	"POSTnGETtrain/internal/userService"
	"POSTnGETtrain/internal/web/users"
	"context"
	"errors"
	"fmt"
)

// UserHandler заготовка для конструктора
type UserHandler struct {
	service userService.UserService // Сервис для работы с пользователями
}

// NewUserHandler создает новый экземпляр UserHandler с заданным сервисом (конструктор)
func NewUserHandler(s userService.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetUsers обрабатывает GET-запрос для получения списка всех пользователей
func (h *UserHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	// Получаем список пользователей из сервиса
	usersList, err := h.service.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Преобразуем пользователей в формат ответа API
	response := make([]users.User, len(usersList))
	for i, u := range usersList {
		response[i] = users.User{
			ID:       u.ID,
			Email:    u.Email,
			Password: u.Password, // Пароль не должен передаваться в API
		}
	}

	// Возвращаем успешный ответ со списком пользователей
	return users.GetUsers200JSONResponse(response), nil
}

// PostUsers обрабатывает POST-запрос для создания нового пользователя
func (h *UserHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	// Проверяем наличие тела запроса
	if request.Body == nil {
		return nil, errors.New("request body is required")
	}

	// Создаем пользователя через сервис
	createdUser, err := h.service.CreateUser(request.Body.Email, request.Body.Password)
	if err != nil {
		if errors.Is(err, userService.ErrEmailExists) {
			return nil, errors.New("email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем успешный ответ с данными созданного пользователя
	return users.PostUsers201JSONResponse{
		ID:       createdUser.ID,
		Email:    createdUser.Email,
		Password: createdUser.Password, // Пароль не должен передаваться в API
	}, nil
}

// PatchUsersId обрабатывает PATCH-запрос для обновления данных пользователя по ID
func (h *UserHandler) PatchUsersId(_ context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	// Обновляем пользователя через сервис
	updatedUser, err := h.service.UpdateUser(
		request.Id,
		request.Body.Email,
		request.Body.Password, // Пароль не должен передаваться в API
	)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Возвращаем успешный ответ с обновленными данными пользователя
	return users.PatchUsersId200JSONResponse{
		ID:       updatedUser.ID,
		Email:    updatedUser.Email,
		Password: updatedUser.Password, // Пароль не должен передаваться в API
	}, nil
}

// DeleteUsersId обрабатывает DELETE-запрос для удаления пользователя по ID
func (h *UserHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	// Удаляем пользователя через сервис
	err := h.service.DeleteUser(request.Id)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	// Возвращаем успешный ответ без содержимого (204 No Content)
	return users.DeleteUsersId204Response{}, nil
}
