package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
)

// GetUserById возвращает текущего пользователя
// @Summary Получить текущего пользователя
// @Description Возвращает профиль авторизованного пользователя
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} User
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Пользователь не найдена"
// @Router /users/me [get]
func GetUserById(w http.ResponseWriter, r *http.Request) {

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

	var user User

	err := db.QueryRow(
		ctx,
		"SELECT id, name, email FROM users WHERE id = $1",
		userID,
	).Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
