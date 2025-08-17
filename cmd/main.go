package main

import (
	"POSTnGETtrain/internal/db"
	"POSTnGETtrain/internal/handlers"
	"POSTnGETtrain/internal/taskService"
	"POSTnGETtrain/internal/web/tasks"
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

	// Инициализация сервисов
	repo := taskService.NewTaskRepository(database)
	service := taskService.NewTaskService(repo)
	handler := handlers.NewHandler(service)

	// Регистрация обработчиков OpenAPI
	strictHandler := tasks.NewStrictHandler(handler, nil)
	tasks.RegisterHandlers(echoServer, strictHandler)

	// Запуск сервера
	err = echoServer.Start("localhost:8080")
	if err != nil {
		log.Fatalf("Could not start: %v", err)
	}
}
