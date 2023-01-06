package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"kafka_test/cache"
	"kafka_test/implementation"
	"kafka_test/kafka"
	"kafka_test/repository"
	"log"
	"net/http"
)

func main() {

	ch := cache.NewCache(0, 0)

	taskRepository := repository.NewTaskRepository()
	taskImplementation := implementation.NewTaskImplementation(taskRepository)

	err := taskImplementation.RestoreCache(ch)
	if err != nil {
		fmt.Println(err)
	}

	kafka.NewConsumer("tasks", taskRepository, ch)

	r := mux.NewRouter()
	r.HandleFunc("/task/{id:[0-9]+}", taskImplementation.GetTask).Methods("GET")
	http.Handle("/", r)
	err = http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
