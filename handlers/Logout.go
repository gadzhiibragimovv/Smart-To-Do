package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
)

// @Description Удаляет все refresh токены пользователя из БД
// @Summary Выйти из системы
// @Description Удаляет все refresh токены...
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

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

	// SQL ЗАПРОС: удаляем токены где user_id = текущий пользователь
	result, err := db.Exec(
		ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)

	if err != nil {
		// Ошибка при удалении из БД
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Сколько токенов удалили (= сколько устройств вышло)
	tokensDeleted := result.RowsAffected()

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправляем успешный ответ клиенту
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Logged out successfully on all devices",
		"tokens_deleted": tokensDeleted, // Количество удаленных токенов
	})
}
