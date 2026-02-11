package utils

import (
	"regexp"
	"strings"
)

// isValidEmail проверяет корректность формата email
// Возвращает true если email валидный, false если нет
func IsValidEmail(email string) bool {
	// Убираем пробелы
	email = strings.TrimSpace(email)

	// Проверяем что email не пустой
	if email == "" {
		return false
	}

	// Регулярное выражение для проверки email
	// Формат: название@домен.зона
	// Примеры валидных: user@mail.ru, test.user@gmail.com, admin+tag@company.co.uk
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

// IsValidPassword проверяет минимальную длину пароля
// Возвращает true если пароль валидный (>= 6 символов), false если нет
func IsValidPassword(password string) bool {
	// Минимум 6 символов
	return len(password) >= 6
}

// IsValidName проверяет корректность имени пользователя
// Возвращает true если имя валидное, false если нет
func IsValidName(name string) bool {
	// Убираем пробелы в начале и конце
	name = strings.TrimSpace(name)

	// Имя не должно быть пустым
	if name == "" {
		return false
	}

	// Имя должно быть от 2 до 50 символов
	if len(name) < 2 || len(name) > 50 {
		return false
	}

	return true
}

// IsValidTaskTitle проверяет корректность названия задачи
// Возвращает true если название валидное, false если нет
func IsValidTaskTitle(title string) bool {
	// Убираем пробелы
	title = strings.TrimSpace(title)

	// Название не должно быть пустым
	if title == "" {
		return false
	}

	// Название должно быть от 1 до 200 символов
	if len(title) > 200 {
		return false
	}

	return true
}

// IsValidTaskDescription проверяет корректность описания задачи
// Возвращает true если описание валидное, false если нет
func IsValidTaskDescription(description string) bool {
	// Описание может быть пустым
	// Но не должно быть слишком длинным (максимум 1000 символов)
	return len(description) <= 1000
}
