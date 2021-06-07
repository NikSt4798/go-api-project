package main

import (
	"etokubernetes/internal/app"
	"etokubernetes/internal/service"
	"flag"
	"log"
)

func main() {
	addr := flag.String("addr", "localhost:80", "specify server lister socket")
	flag.Parse()

	taskService := service.NewTaskService()
	httpApp := app.NewApp(*addr, taskService)
	ch := httpApp.StartHttp()

	err := <-ch

	if err != nil {
		log.Fatalln(err)
	}
}
