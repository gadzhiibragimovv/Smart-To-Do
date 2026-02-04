package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/utils"
	"strings"
)

// LoginHandler авторизует пользователя
// @Summary Войти в систему
// @Description Авторизация по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Email и пароль"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "Неверный JSON"
// @Failure 401 {string} string "Неверные учетные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Ищем пользователя в БД по email
	var user User
	var hashedPassword string

	err = db.QueryRow(
		ctx,
		`SELECT id, name, email, password FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.Id, &user.Name, &user.Email, &hashedPassword)

	if err != nil {
		// Не говорим "user not found", чтобы не раскрывать какие email зарегистрированы
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if err := utils.CheckPassword(hashedPassword, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерируем токены
	accessToken, err := utils.GenerateAccessToken(user.Id)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Не возвращаем пароль клиенту
	user.Password = ""

	// СОЗДАЕМ ОТВЕТ ДЛЯ КЛИЕНТА
	// Заполняем её данными:
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
