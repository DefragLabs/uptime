package main

import (
	"github.com/defraglabs/uptime/internal/api"
	"github.com/defraglabs/uptime/internal/tasks"
)

func main() {
	go uptime.StartScheduler()
	api.StartServer()
}
