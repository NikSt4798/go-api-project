package main

import (
	"bytes"
	"encoding/json"
	"etokubernetes/internal/dto"
	"etokubernetes/internal/task"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	masterAddress := flag.String("master", "http://localhost:80", "specify master address")
	flag.Parse()

	tf := &TaskFetcher{url: *masterAddress + "/task"}
	tr := &TaskReporter{url: *masterAddress + "/task/"}
	for {
		time.Sleep(10 * time.Second)
		t, err := tf.GetTask()
		if err != nil {
			log.Println(err)
			continue
		}

		ts, err := getTaskSolver(t)
		if err != nil {
			t.Status = dto.TaskStatusError
			err = tr.Report(t)
			if err != nil {
				log.Println(err)
			}
			log.Println(err)
			continue
		}

		result, err := ts.Solve(t.Parameters)
		if err != nil {
			log.Println(err)
			continue
		}

		t.Result = result
		t.Status = dto.TaskStatusSuccess

		err = tr.Report(t)
		if err != nil {
			log.Println(err)
		}
	}
}

type TaskFetcher struct {
	url string
}

func (t *TaskFetcher) GetTask() (*dto.Task, error) {
	r, err := http.Get(t.url)
	if err != nil {
		return nil, err
	}

	task := &dto.Task{}
	err = json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

type TaskReporter struct {
	url string
}

func (t *TaskReporter) Report(task *dto.Task) error {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return err
	}
	jsonReader := bytes.NewReader(jsonTask)

	idStr := strconv.Itoa(task.ID)
	_, err = http.Post(t.url+idStr, "application/json", jsonReader)
	if err != nil {
		return err
	}

	return nil
}

func getTaskSolver(t *dto.Task) (task.Task, error) {
	switch t.Type {
	case "simple_numbers":
		return &task.SimpleNumber{}, nil
	default:
		return nil, fmt.Errorf("wrong task type")
	}
}
