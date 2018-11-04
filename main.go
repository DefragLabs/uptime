package main

import (
	"os"

	"github.com/defraglabs/uptime/internal/api"
	"github.com/defraglabs/uptime/internal/tasks"
	log "github.com/sirupsen/logrus"
)

func main() {
	go uptime.StartScheduler()

	setupLogin()
	api.StartServer()
}

func setupLogin() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}
