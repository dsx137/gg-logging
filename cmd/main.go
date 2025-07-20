package main

import (
	"github.com/dsx137/gg-logging/pkg/logging"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.Init()
	logrus.SetLevel(logrus.DebugLevel) // Set the log level to Debug for detailed output

	logrus.Infof("hello world")
}
