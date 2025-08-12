package db

import (
	"POSTnGETtrain/internal/taskService"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Глобальная переменная для связи с БД через GORM
var db *gorm.DB

// Инициализация БД с подключением db к БД
func InitDB() (*gorm.DB, error) {

	// Параметры подключения
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"

	// Подключение к БД
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Автомиграция
	if err = db.AutoMigrate(&taskService.Task{}); err != nil { // Создаем таблицу в БД для модели Task
		log.Fatalf("Could not migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
