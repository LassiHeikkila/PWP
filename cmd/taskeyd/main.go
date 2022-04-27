package main

import (
	"context"
	"fmt"
)

var (
	token = ""
	url   = ""
)

func main() {
	fmt.Println("I don't do anything yet")

	schedule, err := fetchSchedule(token, url)
	if err != nil {
		fmt.Println("error fetching schedule:", err)
		return
	}

	tasks, err := fetchTasks(token, url)
	if err != nil {
		fmt.Println("error fetching tasks:", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executeSchedule(ctx, schedule, tasks)
}
