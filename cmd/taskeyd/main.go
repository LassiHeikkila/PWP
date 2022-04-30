package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

var (
	token = ""
	url   = "https://taskey-service.herokuapp.com"
	org   = "testorg"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	log.Println("fetching schedule")
	schedule, err := fetchSchedule(token, url, org)
	if err != nil {
		fmt.Println("error fetching schedule:", err)
		return
	}
	log.Printf("initial schedule fetched: %#v\n", *schedule)

	log.Println("fetching tasks")
	tasks, err := fetchTasks(token, url, org)
	if err != nil {
		fmt.Println("error fetching tasks:", err)
		return
	}
	log.Println("tasks fetched")
	for name, task := range tasks {
		log.Printf("task %s: %#v content: %#v\n", name, *task, task.Content)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-sc
		cancel()
	}()

	log.Println("executing schedule")
	err = executeSchedule(ctx, schedule, tasks)
	log.Println("executor stopped with error:", err)
}
