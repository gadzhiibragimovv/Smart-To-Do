package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetTaskById возвращает список задач по ID
// @Summary Получить список задач
// @Description Возвращает задачу по её идентификатору
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} Task
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Задача не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks/{id} [get]
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
