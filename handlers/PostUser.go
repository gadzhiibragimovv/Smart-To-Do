package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/utils"
	"strings"
)

// PostUserHandler создает нового пользователя (РЕГИСТРАЦИЯ)
// @Summary Зарегистрировать пользователя
// @Description Создает нового пользователя с хешированным паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "Данные пользователя"
// @Success 201 {object} User "Пользователь успешно создан"
// @Failure 400 {string} string "Неверный JSON или валидация"
// @Failure 409 {string} string "Email уже существует"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/register [post]
func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

	//Тут не нужен JWT middleware, так как пользователь еще не авторизован

	// Создаем структуру для получения данных от клиента
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ВАЛИДАЦИЯ: проверяем, что email не пустой
	// TrimSpace убирает пробелы в начале и конце строки
	if strings.TrimSpace(user.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// ВАЛИДАЦИЯ: проверяем, что имя не пустое
	if strings.TrimSpace(user.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// ХЕШИРУЕМ ПАРОЛЬ перед сохранением в БД
	// Никогда не храним пароли в открытом виде!
	// bcrypt превращает "123456" в "$2a$10$H2m9ZK7P0s2y...Q9x"
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// SQL ЗАПРОС: вставляем нового пользователя в таблицу users
	// $1, $2, $3 - это плейсхолдеры, защищают от SQL-инъекций
	// RETURNING id - PostgreSQL вернет ID созданной записи
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (name, email, password) VALUES ($1,$2,$3) RETURNING id",
		user.Name,
		user.Email,
		hashedPassword,
	).Scan(&user.Id) //Сканируем результат в user_id

	if err != nil {
		// Проверяем, не пытается ли пользователь зарегистрироваться с существующим email
		// PostgreSQL вернет ошибку "duplicate key" если email уже есть в БД
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		// Любая другая ошибка БД
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Очищаем пароль перед отправкой клиенту
	// Клиент НЕ должен получить пароль обратно (даже хешированный)
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
