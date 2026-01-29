package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

// GetTaskById возвращает список задач по ID
// @Summary Получить список задач
// @Description Возвращает задачу текущего пользователя по ID
// @Security BearerAuth
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} Task
// @Failure 400 {string} string "Неверный ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Задача не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks/{id} [get]
func GetTaskById(w http.ResponseWriter, r *http.Request) {

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

	// Достаём userID из context
	// Этот userID туда положил AuthMiddleware после проверки JWT
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		// Если userID нет — значит пользователь не авторизован
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ID задачи из URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	var task Task
	err = db.QueryRow(
		ctx,
		"SELECT id, title, description, created_at, updated_at, user_id FROM tasks WHERE id = $1 and user_id = $2",
		id,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAT,
		&task.UpdatedAT,
		&task.UserID,
	)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
