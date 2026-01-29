package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
)

// GetTaskHandler возвращает список задач
// @Summary Получить список задач
// @Description Возвращает все задачи текущего пользователя
// @Security BearerAuth
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks [get]
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

	// Достаём userID из context (его положил AuthMiddleware)
	// Этот userID туда положил AuthMiddleware после проверки JWT
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		// Если userID нет — значит пользователь не авторизован
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks := []Task{}

	rows, err := db.Query(
		ctx,
		"SELECT id, title, description, created_at, updated_at, user_id FROM tasks WHERE user_id = $1",
		userID,
	)

	if err != nil {
		http.Error(w, "Failed to query tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.CreatedAT,
			&task.UpdatedAT,
			&task.UserID,
		); err != nil {
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
