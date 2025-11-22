package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"smart-todo/handlers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

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

	handlers.SetDB(db)

	r := mux.NewRouter() //Создание роутера

	r.HandleFunc("/tasks", handlers.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", handlers.PostTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.GetTaskById).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	r.HandleFunc("/users", handlers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlers.PostUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.GetUserById).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	fmt.Println("Метод запущен на порту 9090")
	err = http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
