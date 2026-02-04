package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
)

// DeleteUserHandler удаляет текущего пользователя
// @Summary Удалить аккаунт
// @Description Удаляет текущего пользователя и все его задачи (CASCADE)
// @Tags users
// @Security BearerAuth
// @Success 200 {object} map[string]int "ID удалённого пользователя"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Пользователь не найден"
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

	result, err := db.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		userID,
	)

	if err != nil {
		http.Error(w, "Failed to delete user:", http.StatusInternalServerError)
		return
	}

	// ПРОВЕРКА: был ли реально удален пользователь
	// RowsAffected() возвращает количество измененных строк
	// Если 0 - значит задача не найдена или принадлежит другому пользователю
	if result.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"deleted_user_id": userID})
}
