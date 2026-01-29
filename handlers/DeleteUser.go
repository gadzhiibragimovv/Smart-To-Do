package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
)

// DeleteUserHandler удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору
// @Tags users
// @Security BearerAuth
// @Success 200 {integer} int "ID удалённой задачи"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /users/me [delete]
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

	// Достаём userID из context
	// Этот userID туда положил AuthMiddleware после проверки JWT
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		// Если userID нет в context — значит пользователь не авторизован
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		userID,
	)

	if err != nil {
		http.Error(w, "Failed to delete user:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userID)
}
