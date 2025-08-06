package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var task string = "World!" //по умолчанию

// answ - http ответ
func handlerPOST(answ http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Task string `json:"task"`
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(answ, "Bad Request", http.StatusBadRequest)
		return
	}
	task = requestBody.Task

	//подготовка ответа
	response := map[string]string{
		"message": "Task изменен на: " + task, //"message" - ключ к http
	}

	//заворачивание ответа в json
	answ.Header().Set("Content-Type", "application/json")
	json.NewEncoder(answ).Encode(response)
}

func handlerGET(answ http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(answ, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	response := map[string]string{
		"message": "hello, " + task,
	}

	answ.Header().Set("Content-Type", "application/json")
	json.NewEncoder(answ).Encode(response)
}

func main() {
	http.HandleFunc("/", handlerGET)
	http.HandleFunc("/task", handlerPOST)
	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка сервера: ", err)
	}
}
