package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Task struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // "-", чтобы это техническое поле не отображалось в JSON
}

// Соединение с бд
var db *gorm.DB

// Функция инициализации бд
func initDB() {

	// Параметры подключения к "PostgreSQL"
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"

	// Собираем ошибки с соединения
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Автомиграция
	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

// Обработчик GET, выдающий список тасков
func getListTasks(c echo.Context) error {
	//место для задач из БД
	var tasks []Task //<-- сюда

	// Фильтруем неудаленные записи и помещаем таску с помощью Find в tasks
	if err := db.Where("deleted_at IS NULL").Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error."})
	}
	return c.JSON(http.StatusOK, tasks)
}

// Обработчик GET на получение таски по ИД
func getTask(c echo.Context) error {
	//Получаем ID из урла
	id := c.Param("id")

	// Ловим пустой ид
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}

	// переменная типа Таск для записи из БД
	var task Task

	//Поиск по ИД таски из БД, исключая возможность найти удаленные таски
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, task)
}

func postTask(c echo.Context) error {

	//объявление переменной, куда можно положить новую задачу
	var newTask struct {
		Name   string `json:"name"`
		IsDone bool   `json:"is_done"`
	}

	//достаем данные из JSON и кладем их в newTask, заодно проверяем на ошибки
	if err := c.Bind(&newTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	//создаем объект таски для БД
	task := Task{
		ID:     uuid.NewString(), //Генерируем новый ИД
		Name:   newTask.Name,
		IsDone: newTask.IsDone,
	}

	//загружаем в базу данных
	if err := db.Create(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	//сообщаем АПИ, что все сохранилось
	return c.JSON(http.StatusCreated, task)
}

// Удаление задачи
func deleteTask(c echo.Context) error {
	// Получение ID из URL
	id := c.Param("id")

	// Проверка, чтобы не было пустого ID
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Проверяем, не удалена ли уже таска и удаляем её
	result := db.Where("id = ? AND deleted_at IS NULL", id).Delete(&Task{})
	// Ловим ошибки
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	// Было ли изменение в записи?
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task Not Found or deleted"})
	}

	// Уведомляем АПИ, что удаление завершено
	return c.JSON(http.StatusNoContent, nil)
}

func patchTask(c echo.Context) error {
	// Получение ИДа из УРЛ
	id := c.Param("id")

	// Проверка пустого ида
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Ищем задачу по ИДу и загружаем в task(для БД)
	var task Task
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var updates struct {
		Name   *string `json:"name,omitempty"`
		IsDone *bool   `json:"is_done,omitempty"`
	}

	// Перекладываем тело из запроса в структуру updates
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	// Загружаем updates в task с проверкой отсутствующих значений
	if updates.Name != nil {
		task.Name = *updates.Name
	}
	if updates.IsDone != nil {
		task.IsDone = *updates.IsDone
	}

	// Сохраняем изменения в базе
	if err := db.Save(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Отчет об успешном выполнении
	return c.JSON(http.StatusOK, task)
}

func main() {
	initDB()
	e := echo.New()

	e.GET("/tasks", getListTasks)
	e.GET("/tasks/:id", getTask)
	e.POST("/tasks", postTask)
	e.PATCH("/tasks/:id", patchTask)
	e.DELETE("/tasks/:id", deleteTask)

	e.Start("localhost:8080")
}
