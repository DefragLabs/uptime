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
		currentTime := time.Now()
		url := fmt.Sprintf("%s://%s", monitorURL.Protocol, monitorURL.URL)
		start := currentTime

		if int32(currentTime.Minute())%monitorURL.Frequency != 0 {
			continue
		}

		resp, err := http.Get(url)
		duration := time.Since(start)
		if err != nil {
			// Don't fail like this.
			log.Fatal("API ping failed")
		}
		timeStamp := t.Format(time.UnixDate)
		fmt.Println(duration, url, resp.Status, timeStamp)
		db.AddMonitorDetail(monitorURL, resp.Status, timeStamp, duration.String())
	}
}

// StartScheduler runs the scheduler
func StartScheduler() {
	c := make(chan time.Time)
	go func() {
		var frequency = time.Duration(100)
		ticker := time.Tick(frequency)

		for {
			time.Sleep(time.Duration(60 * time.Second))
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
