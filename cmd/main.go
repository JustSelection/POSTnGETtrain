package main

import (
	"POSTnGETtrain/internal/db"
	"POSTnGETtrain/internal/handlers"
	"POSTnGETtrain/internal/taskService"
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

	tskRepo := taskService.NewTaskRepository(database)
	tskService := taskService.NewTaskService(tskRepo)
	tskHandlers := handlers.NewTaskHandler(tskService)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/tasks", tskHandlers.GetListTasks)
	e.GET("/tasks/:id", tskHandlers.GetTask)
	e.POST("/tasks", tskHandlers.PostTask)
	e.PATCH("/tasks/:id", tskHandlers.PatchTask)
	e.DELETE("/tasks/:id", tskHandlers.DeleteTask)

	err = e.Start("localhost:8080")
	if err != nil {
		log.Fatalf("Could not start: %v", err)
	}
}
