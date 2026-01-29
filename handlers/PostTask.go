package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
	"time"
)

// PostTaskHandler создает новую задачу
// @Summary Создать задачу
// @Description Создает новую задачу
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param task body Task true "Данные задачи"
// @Success 201 {object} Task
// @Failure 400 {string} string "Неверный JSON"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks [post]
func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	var task Task

	err := json.NewDecoder(r.Body).Decode(&task) //Из json в Go
	if err != nil {
		// Если userID нет — значит пользователь не авторизован
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task.UserID = userID

	now := time.Now()
	task.CreatedAT = now
	task.UpdatedAT = &now

	// Выполняем SQL-запрос с context
	err = db.QueryRow(
		ctx,
		"INSERT INTO tasks (title, description, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		task.Title,
		task.Description,
		task.UserID,
		task.CreatedAT,
		task.UpdatedAT,
	).Scan(&task.ID)

	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную задачу клиенту
	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(task)                    //Из Go в json
}
