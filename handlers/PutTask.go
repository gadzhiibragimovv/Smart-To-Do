package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	varsUpdate := mux.Vars(r)
	idStrUpdate := varsUpdate["id"]

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStrUpdate)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	now := time.Now()

	var updatedTask Task
	err = db.QueryRow(
		context.Background(),
		"UPDATE tasks SET title = $1, description = $2, updated_at = $3 WHERE id = $4 RETURNING id, title, description, created_at, updated_at",
		task.Title, task.Description, now, id,
	).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.CreatedAT, &updatedTask.UpdatedAT)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask) //Из Go в json
}
