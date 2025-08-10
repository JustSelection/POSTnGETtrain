package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// имитация БД с id:Task
var tasks = make(map[int]Task)

// актуальный ид
var currentID = 1

// Обработчик GET, выдающий список тасков
func getListTasks(c echo.Context) error {

	//место для задач
	taskList := make([]Task, 0, len(tasks))

	//закидываем все таски в taskList
	for _, task := range tasks {
		taskList = append(taskList, task)
	}
	return c.JSON(http.StatusOK, taskList)
}

// Обработчик GET на получение таски по ИД
func getTask(c echo.Context) error {

	//достаем ид из урла, сразу проверяем
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	//проверяем неположительный ид:
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID. Must be a positive"})
	}

	//ищем задачу в "бд"
	task, avail := tasks[id]

	//проверили существование задачи с таким ид
	if !avail {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task Not Found"})
	}

	return c.JSON(http.StatusOK, task)
}

func postTask(c echo.Context) error {

	//объявление переменной, куда можно положить новую задачу
	var newTask Task

	//преобразуем JSON и кладем его содержимое в newTask, заодно проверяем
	if err := c.Bind(&newTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	//присваиваем ид задаче и кладем в "бд"
	newTask.ID = currentID
	tasks[currentID] = newTask

	currentID++ //готов ид для следующей задачи

	return c.JSON(http.StatusCreated, newTask)
}

// удаление задачи
func deleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) //погружаем URL string ID в переменную int id + ошибки
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	//проверка на положительность
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID. Must be a positive"})
	}

	//проверка существования задачи под идом
	_, avail := tasks[id]
	if !avail {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task Not Found"})
	}

	//просто удаление
	delete(tasks, id)
	return c.JSON(http.StatusNoContent, nil)
}

// подрихтовать задачку из списка маленько
func patchTask(c echo.Context) error {

	//грузим в ид ID из урла + проверка на ошибки
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	//ид больше нуля
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID. Must be a positive"})
	}

	//грузим задачу под идом из "бд" в task и проверяем есть ли задача под этим ид
	task, avail := tasks[id]
	if !avail {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task Not Found"})
	}

	//локальная переменная, где будут храниться обновления
	var updates struct {
		Name   string `json:"name,omitempty"` //omitempty - шобы можно было проигнорить пустые значения
		Status string `json:"status,omitempty"`
	}

	//разгружаем JSON в updates
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	//обновляем значения
	if updates.Name != "" {
		task.Name = updates.Name
	}
	if updates.Status != "" {
		task.Status = updates.Status
	}

	//отправляем в бд изменения
	tasks[id] = task

	return c.JSON(http.StatusOK, task)
}

func main() {
	e := echo.New()

	e.GET("/tasks", getListTasks)
	e.GET("/tasks/:id", getTask)
	e.POST("/tasks", postTask)
	e.PATCH("/tasks/:id", patchTask)
	e.DELETE("/tasks/:id", deleteTask)

	e.Start("localhost:8080")
}
