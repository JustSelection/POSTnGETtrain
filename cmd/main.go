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

	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// Инициализация сервисов
	repo := taskService.NewTaskRepository(database)
	service := taskService.NewTaskService(repo)
	handler := handlers.NewHandler(service)

	// Регистрация обработчиков OpenAPI
	strictHandler := tasks.NewStrictHandler(handler, nil)
	tasks.RegisterHandlers(e, strictHandler)

	// Запуск сервера
	err = e.Start("localhost:8080")
	if err != nil {
		log.Fatalf("Could not start: %v", err)
	}
}
