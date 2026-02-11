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

	// Проверка email
	if !utils.IsValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Проверка пароля
	if !utils.IsValidPassword(user.Password) {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Проверка имени(! если не валдиное)
	if !utils.IsValidName(user.Name) {
		http.Error(w, "Name must be between 2 and 50 characters", http.StatusBadRequest)
		return
	}

	// Убираем пробелы
	user.Email = strings.TrimSpace(user.Email)
	user.Name = strings.TrimSpace(user.Name)

	// Проверка дубликата email
	var existingID int
	err = db.QueryRow(ctx,
		`SELECT id FROM users WHERE email = $1`,
		user.Email,
	).Scan(&existingID)

	if err == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
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
