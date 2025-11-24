package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID          int        `json:"id"`          //ID - Уникальный индефикатор
	Title       string     `json:"title"`       //Title - Заголовок задачи
	Description string     `json:"description"` //Description - Краткое описание задачи
	CreatedAT   time.Time  `json:"created_at"`  //CreatedAT - Дата и время создания задачи
	UpdatedAT   *time.Time `json:"updated_at"`  //UpdatedAT - Дата и время последнего обновления задачи
}

var db *pgx.Conn

func SetDB(database *pgx.Conn) {
	db = database
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{}

	rows, err := db.Query(context.Background(), "SELECT id, title, description, created_at, updated_at FROM tasks")
	if err != nil {
		http.Error(w, "Failed to query tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAT, &task.UpdatedAT); err != nil {
			http.Error(w, "Error scanning user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
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
	err = db.QueryRow(
		context.Background(),
		"SELECT id, title, description, created_at, updated_at FROM tasks WHERE id = $1",
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAT, &task.UpdatedAT)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task) //Из json в Go
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	now := time.Now()
	task.CreatedAT = now
	task.UpdatedAT = &now

	err = db.QueryRow(
		context.Background(),
		"INSERT INTO tasks (title, description, created_at, updated_at) VALUES ($1,$2,$3,$4) RETURNING id",
		task.Title,
		task.Description,
		task.CreatedAT,
		task.UpdatedAT,
	).Scan(&task.ID)

	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

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

	_, err = db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed from delete tasks:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(id)                      //Из Go в json
}
