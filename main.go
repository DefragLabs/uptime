package main

import (
	"github.com/dineshs91/uptime/internal/api"
	"github.com/dineshs91/uptime/internal/tasks"
)

func main() {
	api.StartServer()
	uptime.PingServer()
}
