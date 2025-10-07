package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Task struct {
	ID       int    `json:"id"`       //ID - Уникальный индефикатор
	UserID   int    `json:"userid"`   //UserID - Связь с пользователем, пользователь может видеть/изменять свои данные
	Title    string `json:"title"`    //Title - заголовок задачи
	IsDone   bool   `json:"isdone"`   //IsDone - Выполнение задачи(true-выполнено/false-ложно)
	Priority int    `json:"priority"` //Priority - Приоритет - (0-неуказан,1-низкий,2-средний,3-высокий)
	Category string `json:"category"` //Category - Категории(дом,работа)

}

var Tasks = []Task{
	{ID: 1, UserID: 1, Title: "Задача 1", IsDone: false, Priority: 0, Category: "Работа"},
	{ID: 2, UserID: 2, Title: "Задача 2", IsDone: true, Priority: 1, Category: "Дом"},
	{ID: 3, UserID: 3, Title: "Задача 3", IsDone: false, Priority: 2, Category: "Работа"},
	{ID: 4, UserID: 4, Title: "Задача 4", IsDone: true, Priority: 3, Category: "Дом"},
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(Tasks)                   //Из Go в json
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task) //Из json в Go
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	Tasks = append(Tasks, task)
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(task)                    //Из Go в json
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	varsDelete := mux.Vars(r)
	idStrDelete := varsDelete["id"]

	id, err := strconv.Atoi(idStrDelete)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	for i, task := range Tasks {
		if task.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
			json.NewEncoder(w).Encode(Tasks)                   //Из Go в json
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "error missing id in path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	for _, task := range Tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func main() {
	err := godotenv.Load() //загрузка .env файла
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter() //Создание роутера

	r.HandleFunc("/tasks", GetTaskHandler).Methods("Get")
	r.HandleFunc("/tasks", PostTaskHandler).Methods("Post")
	r.HandleFunc("/tasks/{id}", GetTaskById).Methods("Get")
	r.HandleFunc("/tasks/{id}", DeleteTaskHandler).Methods("Delete")
	fmt.Println("Метод запущен на порту 9090")
	err = http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
