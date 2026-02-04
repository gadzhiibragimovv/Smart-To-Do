package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "smart-todo/docs" //Без этого сваггер не работает
	"smart-todo/handlers"
	"smart-todo/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Smart To-Do API
// @version 1.0
// @description API для управления задачами и пользователями
// @host localhost:9090
// @BasePath /
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

	// Публичные маршруты (БЕЗ JWT)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/auth/register", handlers.PostUserHandler).Methods("POST")
	r.HandleFunc("/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", handlers.RefreshTokenHandler).Methods("POST")

	// Создаём подроутер (группу маршрутов)
	// PathPrefix("/") означает: все пути, начинающиеся с "/"
	auth := r.PathPrefix("/").Subrouter()

	// Подключаем middleware ТОЛЬКО к этому подроутеру
	// Каждый запрос к маршрутам ниже сначала пройдёт AuthMiddleware
	auth.Use(middleware.AuthMiddleware)

	// GET /tasks — получить все задачи текущего пользователя
	auth.HandleFunc("/tasks", handlers.GetTaskHandler).Methods("GET")
	// POST /tasks — создать задачу для текущего пользователя
	auth.HandleFunc("/tasks", handlers.PostTaskHandler).Methods("POST")
	// GET /tasks/{id} — получить задачу по ID (только свою)
	auth.HandleFunc("/tasks/{id}", handlers.GetTaskById).Methods("GET")
	// PUT /tasks/{id} — обновить задачу (только свою)
	auth.HandleFunc("/tasks/{id}", handlers.PutTaskHandler).Methods("PUT")
	// DELETE /tasks/{id} — удалить задачу (только свою)
	auth.HandleFunc("/tasks/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	// GET /users/me — получить профиль текущего пользователя
	auth.HandleFunc("/users/me", handlers.GetUserById).Methods("GET")
	// PUT /users/me — обновить профиль текущего пользователя
	auth.HandleFunc("/users/me", handlers.PutUserHandler).Methods("PUT")
	// DELETE /users/me — удалить текущего пользователя
	auth.HandleFunc("/users/me", handlers.DeleteUserHandler).Methods("DELETE")

	fmt.Println("Метод запущен на порту 9090")
	err = http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
