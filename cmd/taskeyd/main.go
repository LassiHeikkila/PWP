package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	config = Config{}
)

func main() {
	log.SetFlags(log.Ldate | log.LUTC | log.Lshortfile | log.Lmicroseconds)
	conf := flag.String("c", "", "path to configuration JSON file")
	demoMode := flag.Bool("demo", false, "go slow for demo purposes")
	flag.Parse()

	if *conf == "" {
		log.Println("you must provide configuration file")
		flag.Usage()
		return
	}

	if err := loadConfig(*conf, &config); err != nil {
		log.Println("error loading configuration:", err)
	}

	log.Println("checking token validity...")

	if err := checkToken(config.AccessToken, config.URL, config.Organization); err != nil {
		log.Println("error checking token validity:", err)
		return
	}
	if *demoMode {
		time.Sleep(time.Second)
		log.Println("token check passed!")
		time.Sleep(2 * time.Second)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	log.Println("fetching schedule")
	schedule, err := fetchSchedule(config.AccessToken, config.URL, config.Organization)
	if err != nil {
		log.Println("error fetching schedule:", err)
		return
	}
	if *demoMode {
		time.Sleep(time.Second)
		log.Println("schedule fetched:")
		sb, _ := json.MarshalIndent(schedule, "", "  ")
		_, _ = log.Writer().Write(sb)
		_, _ = log.Writer().Write([]byte("\n"))
		time.Sleep(2 * time.Second)
	}

	log.Println("fetching tasks")
	tasks, err := fetchTasks(config.AccessToken, config.URL, config.Organization)
	if err != nil {
		log.Println("error fetching tasks:", err)
		return
	}
	if *demoMode {
		time.Sleep(time.Second)
		log.Println("tasks fetched:")
		tb, _ := json.MarshalIndent(tasks, "", "  ")
		_, _ = log.Writer().Write(tb)
		_, _ = log.Writer().Write([]byte("\n"))
		time.Sleep(2 * time.Second)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-sc
		cancel()
	}()

	if *demoMode {
		time.Sleep(time.Second)
		log.Println("executing schedule...")
		time.Sleep(2 * time.Second)
	}
	err = executeSchedule(ctx, schedule, tasks)
	log.Println("executor stopped with error:", err)
}
