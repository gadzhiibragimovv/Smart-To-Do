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

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetTaskHandler(w, r)
	} else if r.Method == http.MethodPost {
		PostTaskHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(Tasks)                   //Из Go в json
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task) //Из json в Go
	Tasks = append(Tasks, task)
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(task)                    //Из Go в json
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

	var task Task
	for _, t := range Tasks {
		if t.ID == id {
			task = t
			return
		}
	}
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(task)                    //Из Go в json
}

func main() {
	err := godotenv.Load() //загрузка .env файла
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter() //Создание роутера

	http.HandleFunc("\tasks", TaskHandler)
	fmt.Println("Метод запущен на порту 9090")
	http.ListenAndServe("9090", r)
}
