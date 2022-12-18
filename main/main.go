package main

import (
	"fmt"
	"kafka_test/postgres"
	"log"
	"time"
)

func main() {
	taskRepository := postgres.NewTaskRepository()

	//kafka.NewConsumer("test")

	var t postgres.DbTask
	t.Task = "Test"
	t.Time = time.Now()

	id, err := taskRepository.PostTaskToDB(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

}
