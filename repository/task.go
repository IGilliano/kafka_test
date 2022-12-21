package repository

import (
	"database/sql"
	"fmt"
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
	fmt.Println("Got to DB")
	var task []*DbTask
	rows, err := tr.db.Query("SELECT * FROM tasks WHERE id = $1", id)

	fmt.Println(err)

	for rows.Next() {
		fmt.Println("Even got here")
		var taskScan DbTask
		err = rows.Scan(&taskScan.Task, &taskScan.Time, &taskScan.Id)
		if err == nil {
			task = append(task, &taskScan)
		}

	}
	fmt.Println(err)
	return task, err
}
