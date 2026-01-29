package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"smart-todo/middleware"

	"github.com/jackc/pgx/v5"
)

// PutUserHandler обновляет текущего пользователя
// @Summary Обновить текущего пользователя
// @Description Обновляет профиль авторизованного пользователя
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param task body User true "Данные пользователей"
// @Success 200 {object} User
// @Failure 400 {string} string "Неверный JSON или ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /users/me [put]
func PutUserHandler(w http.ResponseWriter, r *http.Request) {
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

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var updatedUser User

	err = db.QueryRow(
		ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email",
		user.Name,
		user.Email,
		userID,
	).Scan(&updatedUser.Id,
		&updatedUser.Name,
		&updatedUser.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser) //Из Go в json
}
