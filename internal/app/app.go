package app

import (
	"encoding/json"
	"etokubernetes/internal/dto"
	"etokubernetes/internal/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"strings"
	"bytes"
)

type App interface {
	StartHttp() chan error
	// TODO StopHttp graceful shutdown
}

type app struct {
	router      *mux.Router
	addr        string
	tasks       sync.Map
	taskService service.TaskService
}

func NewApp(addr string, taskService service.TaskService) App {
	app := &app{
		router:      mux.NewRouter(),
		addr:        addr,
		taskService: taskService,
	}

	app.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})

	app.router.Path("/task").Methods(http.MethodPost).HandlerFunc(app.createTask)
	app.router.Path("/task/{id:[\\d+]}").Methods(http.MethodPost).HandlerFunc(app.reportTask)
	app.router.Path("/task").Methods(http.MethodGet).HandlerFunc(app.getTask)
	app.router.Path("/table").Methods(http.MethodGet).HandlerFunc(app.listTasks)

	return app
}

func (a *app) StartHttp() chan error {
	ch := make(chan error)

	go func() {
		ch <- http.ListenAndServe(a.addr, a.router)
	}()

	log.Printf("http server started on %s", a.addr)

	return ch
}

func (a *app) getTask(w http.ResponseWriter, r *http.Request) {
	task, err := a.taskService.GetTask(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (a *app) createTask(w http.ResponseWriter, r *http.Request) {
	task := &dto.Task{}

	json.NewDecoder(r.Body).Decode(task)

	task, err := a.taskService.CreateTask(r.Context(), task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Errorf("task create error")
		return
	}

	_ = json.NewEncoder(w).Encode(task)
}

func (a *app) reportTask(w http.ResponseWriter, r *http.Request) {
	task := &dto.Task{}

	json.NewDecoder(r.Body).Decode(task)

	err := a.taskService.ReportTask(r.Context(), task)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (a *app) listTasks(w http.ResponseWriter, r *http.Request) {
	tableBuildr := strings.Builder{}

	tableBuildr.WriteString("<!DOCTYPE html><html><head><meta charset=\"UTF-8\"></head><body><table><tr><th>Type</th><th>Status</th><th>Parameters</th><th>Result</th></tr>")

	tasks, err := a.taskService.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, v := range tasks {
		pBytes := make([]byte, 100)
		rBytes := make([]byte, 100)
		p := bytes.NewBuffer(pBytes)
		rString := bytes.NewBuffer(rBytes)

		_ = json.NewEncoder(p).Encode(v.Parameters)
		_ = json.NewEncoder(rString).Encode(v.Result)
		tableBuildr.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			v.Type, dto.TaskStatuses[v.Status], p.String(), rString.String()))
	}

	tableBuildr.WriteString("</table></body></html>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(tableBuildr.String()))
}
