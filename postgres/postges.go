package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DbTask struct {
	Id   int       `json:"id"`
	Task string    `json:"task"`
	Time time.Time `json:"time"`
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository() TaskRepository {
	db, err := sql.Open("pgx", "postgres://postgres:456123789@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return TaskRepository{db: db}
}

func (tr *TaskRepository) PostTaskToDB(task string, time time.Time) (int, error) {
	var id int
	err := tr.db.QueryRow("INSERT INTO tasks (task, time) VALUES ($1, $2) RETURNING id", task, time).Scan(&id)
	return id, err
}
