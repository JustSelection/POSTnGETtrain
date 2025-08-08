package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

var tasks = make(map[int]Task)
var currentID = 1

func handlerMethod(answ http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		handlerPOST(answ, req)
	case http.MethodPatch:
		handlerPATCH(answ, req)
	case http.MethodDelete:
		handlerDELETE(answ, req)
	default:
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handlerPOST(answ http.ResponseWriter, req *http.Request) {

	//проверка совпадения метода
	if req.Method != http.MethodPost {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//декодирование JSON
	var newTask Task                                  //создание переменной для задачи на основе структуры Task
	err := json.NewDecoder(req.Body).Decode(&newTask) //создание переменной, которая собирает ошибки
	if err != nil {                                   //при декодировании
		http.Error(answ, "Bad Request", http.StatusBadRequest)
		return
	}

	//присвоение ID
	newTask.ID = currentID
	tasks[currentID] = newTask //добавление таски с айди в глобальную мапу "базу данных"
	currentID++

	//создание ответа
	response := map[string]interface{}{ //использую значение interface{}, чтобы не возникло конфликтов типов
		"id":     newTask.ID,
		"text":   newTask.Text,
		"status": newTask.Status,
	}

	//завернуть в JSON и отправить

	answ.Header().Set("Content-Type", "application/json")
	answ.WriteHeader(http.StatusCreated)
	json.NewEncoder(answ).Encode(response)
}

func handlerGET(answ http.ResponseWriter, req *http.Request) {

	//проверяем совместимость методов
	if req.Method != http.MethodGet {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//получение ID из URL
	idStr := req.URL.Query().Get("id")

	//проверка пустого ID?
	if idStr == "" {
		taskList := make([]Task, 0, len(tasks))
		for _, task := range tasks {
			taskList = append(taskList, task)
		}

		answ.Header().Set("Content-Type", "application/json")
		json.NewEncoder(answ).Encode(taskList)
		return
	}

	//преобразуем полученный string-ID в int-ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(answ, "Invalid ID", http.StatusBadRequest)
		return
	}

	//поиск задачи по списку tasks
	task, avail := tasks[id] // avail - показывает наличие задачи(bool)
	if !avail {
		http.Error(answ, "Task Not Found", http.StatusNotFound)
	}

	//заворачиваем в JSON и отправляем задачу клиенту
	answ.Header().Set("Content-Type", "application/json")
	json.NewEncoder(answ).Encode(task)
}

func handlerPATCH(answ http.ResponseWriter, req *http.Request) {

	//проверка метода на вшивость
	if req.Method != http.MethodPatch {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//получаем ID задачи
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		http.Error(answ, "Cannot be empty", http.StatusBadRequest)
		return
	}

	//преобразую ID-string to ID-int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(answ, "Invalid ID", http.StatusBadRequest)
		return
	}

	//поиск задачи | есть/нет задача = avail(bool)
	task, avail := tasks[id]
	if !avail {
		http.Error(answ, "Task Not Found", http.StatusNotFound)
		return
	}

	//обработка и декодирование изменений в тасках
	var updates struct {
		Text   string `json:"text,omitempty"` // omitempty - позволяет пропускать пустые/не передаваемые зн-я
		Status string `json:"status,omitempty"`
	}
	err = json.NewDecoder(req.Body).Decode(&updates)

	//применяем обновления
	if updates.Text != "" {
		task.Text = updates.Text
	}
	if updates.Status != "" {
		task.Status = updates.Status
	}

	//запись в "БД"
	tasks[id] = task

	//формируем и отправляем ответ
	response := map[string]interface{}{
		"id":     id,
		"text":   task.Text,
		"status": task.Status,
	}
	answ.Header().Set("Content-Type", "application/json")
	json.NewEncoder(answ).Encode(response)

}

func handlerDELETE(answ http.ResponseWriter, req *http.Request) {

	//проверка метода
	if req.Method != http.MethodDelete {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//прием ID
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		http.Error(answ, "Cannot be empty", http.StatusBadRequest)
		return
	}

	//преобразование id string - id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(answ, "Invalid ID", http.StatusBadRequest)
		return
	}

	//найти задачу
	_, avail := tasks[id] //'_' - потому что нет необходимости использовать таску
	if !avail {
		http.Error(answ, "Task Not Fount", http.StatusNotFound)
		return
	}

	//удоли задачу!
	delete(tasks, id)
}

func main() {
	http.HandleFunc("/", handlerGET)
	http.HandleFunc("/tasks", handlerMethod)
	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка сервера: ", err)
	}
}
