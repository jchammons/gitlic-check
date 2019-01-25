package main

import (
	"github.com/sirupsen/logrus"

	"github.com/solarwinds/gitlic-check/cmd"
)

var log *logrus.Logger

func main() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	cmd.RegisterLogger(log)
	cmd.Execute()
}
