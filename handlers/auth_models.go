package handlers

//Структуры для авторизации

// LoginRequest - запрос на вход в систему
type LoginRequest struct {
	Email    string `json:"email"`    // Email пользователя
	Password string `json:"password"` // Пароль (не хешированный)
}

// LoginResponse - ответ при успешном входе
type LoginResponse struct {
	AccessToken  string `json:"access_token"`  // JWT токен для доступа (24 часа)
	RefreshToken string `json:"refresh_token"` // JWT токен для обновления (7 дней)
	User         User   `json:"user"`          //Информация о пользователе из файла user.go
}

// RefreshRequest - запрос на обновление access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"` //Refresh токен
}

// RefreshResponse - ответ с новым access token
type RefreshResponse struct {
	AccessToken string `json:"access_token"` //Новый JWT токен
}
