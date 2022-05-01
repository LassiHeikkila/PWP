package main

import (
	"context"
	"encoding/json"
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
	log.SetFlags(log.Ldate | log.LUTC | log.Lshortfile | log.Lmicroseconds)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	log.Println("fetching schedule")
	schedule, err := fetchSchedule(token, url, org)
	if err != nil {
		fmt.Println("error fetching schedule:", err)
		return
	}

	log.Println("fetching tasks")
	tasks, err := fetchTasks(token, url, org)
	if err != nil {
		fmt.Println("error fetching tasks:", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-sc
		cancel()
	}()

	log.Println("executing schedule:")
	sb, _ := json.MarshalIndent(schedule, "", "  ")
	_, _ = log.Writer().Write(sb)
	_, _ = log.Writer().Write([]byte("\n"))

	log.Println("defined tasks:")
	tb, _ := json.MarshalIndent(tasks, "", "  ")
	_, _ = log.Writer().Write(tb)
	_, _ = log.Writer().Write([]byte("\n"))

	err = executeSchedule(ctx, schedule, tasks)
	log.Println("executor stopped with error:", err)
}
