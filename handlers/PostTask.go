package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

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
