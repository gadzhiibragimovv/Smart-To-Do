package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Task struct {
	ID       int    `json:"id"`       //ID - Уникальный индефикатор
	UserID   int    `json:"userid"`   //UserID - Связь с пользователем, пользователь может видеть/изменять свои данные
	Title    string `json:"title"`    //Title - заголовок задачи
	IsDone   bool   `json:"isdone"`   //IsDone - Выполнение задачи(true-выполнено/false-не выполнено)
	Priority int    `json:"priority"` //Priority - Приоритет - (0-неуказан,1-низкий,2-средний,3-высокий)
	Category string `json:"category"` //Category - Категории(дом,работа)
}

var Tasks = []Task{
	{ID: 1, UserID: 1, Title: "Задача 1", IsDone: false, Priority: 0, Category: "Работа"},
	{ID: 2, UserID: 2, Title: "Задача 2", IsDone: true, Priority: 1, Category: "Дом"},
	{ID: 3, UserID: 3, Title: "Задача 3", IsDone: false, Priority: 2, Category: "Работа"},
	{ID: 4, UserID: 4, Title: "Задача 4", IsDone: true, Priority: 3, Category: "Дом"},
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"users"`
}

var db *pgx.Conn

var Users = []User{
	{Id: 1, Name: "Пользователь 1"},
	{Id: 2, Name: "Пользователь 2"},
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(Tasks)                   //Из Go в json
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "error missing in path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow(
		context.Background(),
		"SELECT id, name FROM users WHERE id = $1",
		id,
	).Scan(&user.Id, &user.Name)

	for _, user := range Users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	Users = append(Users, user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	varsDelete := mux.Vars(r)
	idStrDelete := varsDelete["id"]

	id, err := strconv.Atoi(idStrDelete)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	for i, user := range Users {
		if user.Id == id {
			Users = append(Users[:i], Users[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Users)
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

	//Подключение к БД
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	var insertedID int
	err = db.QueryRow(
		context.Background(),
		`INSERT INTO users (name) VALUES ($1) 
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
         RETURNING id`,
	).Scan(&insertedID)

	r := mux.NewRouter() //Создание роутера

	r.HandleFunc("/tasks", GetTaskHandler).Methods("Get")
	r.HandleFunc("/tasks", PostTaskHandler).Methods("Post")
	r.HandleFunc("/tasks/{id}", GetTaskById).Methods("Get")
	r.HandleFunc("/tasks/{id}", DeleteTaskHandler).Methods("Delete")

	r.HandleFunc("/users", GetUserHandler).Methods("Get")
	r.HandleFunc("/users", PostUserHandler).Methods("Post")
	r.HandleFunc("/users/{id}", GetUserById).Methods("Get")
	r.HandleFunc("/users/{id}", DeleteUserHandler).Methods("Delete")

	fmt.Println("Метод запущен на порту 9090")
	err = http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
