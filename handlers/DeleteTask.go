package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteTaskHandler удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {integer} int "ID удалённой задачи"
// @Failure 400 {string} string "Неверный ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

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

	varsDelete := mux.Vars(r)
	idStrDelete := varsDelete["id"]

	id, err := strconv.Atoi(idStrDelete)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		ctx,
		"DELETE FROM tasks WHERE id = $1 and user_id = $2",
		id,
		userID,
	)

	if err != nil {
		http.Error(w, "Failed from delete tasks:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(id)                      //Из Go в json
}
