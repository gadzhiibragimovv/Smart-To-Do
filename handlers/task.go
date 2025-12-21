package handlers

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID          int        `json:"id"`          //ID - Уникальный индефикатор
	Title       string     `json:"title"`       //Title - Заголовок задачи
	Description string     `json:"description"` //Description - Краткое описание задачи
	CreatedAT   time.Time  `json:"created_at"`  //CreatedAT - Дата и время создания задачи
	UpdatedAT   *time.Time `json:"updated_at"`  //UpdatedAT - Дата и время последнего обновления задачи
}

var db *pgx.Conn

func SetDB(database *pgx.Conn) {
	db = database
}
