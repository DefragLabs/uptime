package uptime

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dineshs91/uptime/internal/db"
)

func pingURL(t time.Time) {
	monitoringURLS := db.GetMonitoringURLS()

	for _, monitorURL := range monitoringURLS {
		url := fmt.Sprintf("%s://%s", monitorURL.Protocol, monitorURL.URL)
		resp, err := http.Get(url)
		if err != nil {
			// Don't fail like this.
			log.Fatal("API ping failed")
		}
		fmt.Println(url, resp.Status, t.Format(time.UnixDate))
	}
}

// StartScheduler runs the scheduler
func StartScheduler() {
	c := make(chan time.Time)
	go func() {
		var frequency = time.Duration(100)
		ticker := time.Tick(frequency)

		for {
			time.Sleep(time.Duration(180 * time.Second))
			c <- <-ticker
		}
	}()

	for {
		select {
		case t := <-c:
			pingURL(t)
		case <-time.After(time.Duration(300 * time.Second)):
			// This case acts as a timeout.
			fmt.Println("Ending")
		}
	}
}
