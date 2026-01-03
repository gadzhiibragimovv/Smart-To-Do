package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetTaskHandler возвращает список задач
// @Summary Получить список задач
// @Description Возвращает все задачи
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks [get]

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
