package repository

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

type ITaskRepository interface {
	PostTaskToDB(string, time.Time) (int, error)
	GetTaskFromDB(id int) ([]*DbTask, error)
	GetTasks() ([]*DbTask, error)
}

func NewTaskRepository() ITaskRepository {
	db, err := sql.Open("pgx", "postgres://postgres:456123789@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) PostTaskToDB(task string, time time.Time) (int, error) {
	var id int
	err := tr.db.QueryRow("INSERT INTO tasks (task, time) VALUES ($1, $2) RETURNING id", task, time).Scan(&id)
	return id, err
}

func (tr *TaskRepository) GetTaskFromDB(id int) ([]*DbTask, error) {
	var task []*DbTask
	rows, err := tr.db.Query("SELECT * FROM tasks WHERE id = $1", id)

	for rows.Next() {
		var taskScan DbTask
		err = rows.Scan(&taskScan.Task, &taskScan.Time, &taskScan.Id)
		if err == nil {
			task = append(task, &taskScan)
		}

	}
	return task, err
}

func (tr *TaskRepository) GetTasks() ([]*DbTask, error) {
	var tasks []*DbTask
	rows, err := tr.db.Query("SELECT * FROM tasks")

	for rows.Next() {
		var taskScan DbTask
		err = rows.Scan(&taskScan.Task, &taskScan.Time, &taskScan.Id)
		if err == nil {
			tasks = append(tasks, &taskScan)
		}
	}
	return tasks, err
}
