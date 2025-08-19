package main

import (
	"POSTnGETtrain/internal/db"
	"POSTnGETtrain/internal/handlers"
	"POSTnGETtrain/internal/taskService"
	"POSTnGETtrain/internal/userService"
	"POSTnGETtrain/internal/web/tasks"
	"POSTnGETtrain/internal/web/users"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	echoServer := echo.New()

	// Middleware
	echoServer.Use(middleware.CORS())
	echoServer.Use(middleware.Logger())

	// Инициализация сервисов задач
	tskRepo := taskService.NewTaskRepository(database)
	tskService := taskService.NewTaskService(tskRepo)
	tskHandler := handlers.NewHandler(tskService)

	// Инициализация сервисов пользователей
	usrRepo := userService.NewUserRepository(database)
	usrService := userService.NewUserService(usrRepo)
	usrHandler := handlers.NewUserHandler(usrService)

	// Регистрация обработчиков OpenAPI
	taskStrictHandler := tasks.NewStrictHandler(tskHandler, nil)
	tasks.RegisterHandlers(echoServer, taskStrictHandler)

	userStrictHandler := users.NewStrictHandler(usrHandler, nil)
	users.RegisterHandlers(echoServer, userStrictHandler)
	// Запуск сервера
	err = echoServer.Start("localhost:8080")
	if err != nil {
		log.Fatalf("Could not start: %v", err)
	}
}
