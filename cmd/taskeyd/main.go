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
	url   = ""
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	log.Println("fetching schedule")
	schedule, err := fetchSchedule(token, url)
	if err != nil {
		fmt.Println("error fetching schedule:", err)
		return
	}
	log.Println("initial schedule fetched")

	log.Println("fetching tasks")
	tasks, err := fetchTasks(token, url)
	if err != nil {
		fmt.Println("error fetching tasks:", err)
		return
	}
	log.Println("tasks fetched")

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
