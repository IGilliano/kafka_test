package implementation

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"kafka_test/repository"
	"log"
	"net/http"
	"strconv"
)

type TaskImplementation struct {
	taskRepository repository.ITaskRepository
}

type ITaskImplementation interface {
	GetTask(http.ResponseWriter, *http.Request)
}

func NewTaskImplementation(tr repository.ITaskRepository) ITaskImplementation {
	return &TaskImplementation{taskRepository: tr}
}

func getId(r *http.Request) int {
	vars := mux.Vars(r)
	stringid := vars["id"]
	id, err := strconv.Atoi(stringid)
	if err != nil {
		log.Fatal(err)
	}
	return id

}

func (ti TaskImplementation) GetTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get request!")
	id := getId(r)
	task, err := ti.taskRepository.GetTaskFromDB(id)
	taskByte, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, newErr := w.Write([]byte(err.Error()))
		if newErr != nil {
			log.Fatal(newErr)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write(taskByte)
	if errWrite != nil {
		log.Fatal(errWrite)
	}

}
