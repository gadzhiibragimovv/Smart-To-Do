package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/middleware"
	"smart-todo/utils"
)

// @Summary Выйти из системы
// @Description Удаляет конкретный refresh token из БД (выход на одном устройстве)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param logout body LogoutRequest true "Refresh token"
// @Success 200 {object} map[string]string "Успешный выход"
// @Failure 400 {string} string "Неверный JSON"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Token does not belong to this user"
// @Failure 404 {string} string "Token not found"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Закрываем тело запроса после выполнения функции

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

	// Создаем переменную для получения refresh token от клиента
	var req LogoutRequest

	// Декодируем JSON из тела запроса
	// Клиент отправил: {"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ВАЛИДАЦИЯ: проверяем, что refresh token не пустой
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// ПРОВЕРЯЕМ REFRESH TOKEN - извлекаем userID из токена
	// ParseRefreshToken делает:
	// 1. Проверяет подпись токена
	// 2. Проверяет срок действия
	// 3. Извлекает userID
	userIDFromToken, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		// Токен невалидный или истек
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// БЕЗОПАСНОСТЬ: проверяем, что user_id из access token совпадает с user_id из refresh token
	// Это гарантирует, что пользователь удаляет СВОЙ токен, а не чужой
	// Пример атаки без этой проверки:
	// 1. Хакер крадет чужой refresh token
	// 2. Хакер использует свой access token + чужой refresh token
	// 3. Хакер удаляет чужой refresh token через logout
	// 4. Жертва больше не может обновить токен
	if userID != userIDFromToken {
		// user_id не совпадают - пользователь пытается удалить чужой токен!
		http.Error(w, "Token does not belong to this user", http.StatusForbidden)
		return
	}

	// SQL ЗАПРОС: удаляем КОНКРЕТНЫЙ токен
	result, err := db.Exec(
		ctx,
		`DELETE FROM refresh_tokens WHERE token = $1 AND user_id = $2`,
		req.RefreshToken, // $1 - удаляем ЭТОТ токен
		userID,           // $2 - если принадлежит этому пользователю
	)

	if err != nil {
		// Ошибка при удалении из БД
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// ПРОВЕРЯЕМ: был ли реально удален токен
	// RowsAffected() возвращает количество удаленных строк
	// Если 0 - значит токен не найден в БД
	if result.RowsAffected() == 0 {
		// Токен не найден в БД
		// Возможные причины:
		// 1. Токен уже был удален раньше (повторный logout)
		// 2. Токен никогда не был сохранен (поддельный)
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправляем успешный ответ клиенту
	// Клиент получит: {"message": "Logged out successfully"}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
