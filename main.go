package main

import (
	"os"

	"github.com/defraglabs/uptime/internal/api"
	"github.com/defraglabs/uptime/internal/tasks"
	log "github.com/sirupsen/logrus"
)

func init() {
	setupLogin()
}

func main() {
	go uptime.StartScheduler()

	api.StartServer()
}

func setupLogin() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
