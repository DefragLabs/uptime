package uptime

import (
	"fmt"
	"time"
)

// PingServer This function pings any server, and stores the status code.
func PingServer() {
	var frequence = time.Duration(10)
	ticker := time.NewTicker(frequence)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Printing")
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}
