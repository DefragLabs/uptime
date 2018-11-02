package main

import (
	"github.com/defraglabs/uptime/internal/api"
)

func main() {
	go uptime.StartScheduler()
	api.StartServer()
}
