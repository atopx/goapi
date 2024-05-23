package worker

import (
	"fmt"
	"goapi/conf"
	"goapi/internal/worker/health"
	"log"

	"github.com/robfig/cron/v3"
)

func Start(cfg map[string]*conf.WorkerConfig, errors chan<- error) {
	crontab := cron.New(cron.WithSeconds())

	for key, worker := range cfg {
		if worker.Enable {
			var job func()
			switch key {
			case "health":
				job = health.Health
			default:
				log.Printf("ignore unknown worker: %s", key)
				continue
			}
			_, err := crontab.AddFunc(worker.Spec, job)
			if err != nil {
				errors <- fmt.Errorf("registe worker %s failed with %s", key, err)
				break
			}
			log.Printf("register worker: %s", key)
		} else {
			log.Printf("ignore disabled worker: %s", key)
		}
	}

	go func() {
		log.Println("start workers")
		crontab.Run()
	}()
}
