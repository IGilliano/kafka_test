package main

import (
	"github.com/gorilla/mux"
	"kafka_test/implementation"
	"kafka_test/kafka"
	"kafka_test/repository"
	"log"
	"net/http"
)

func main() {

	kafka.NewConsumer("tasks")

	taskRepository := repository.NewTaskRepository()
	taskImplementation := implementation.NewTaskImplementation(taskRepository)
	r := mux.NewRouter()
	r.HandleFunc("/task/{id:[0-9]+}", taskImplementation.GetTask).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
