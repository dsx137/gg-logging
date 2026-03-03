package main

import (
	"github.com/dsx137/gg-logging/pkg/gglogging"
	"github.com/sirupsen/logrus"
)

func main() {
	gglogging.Init()
	logrus.SetLevel(logrus.DebugLevel) // Set the log level to Debug for detailed output

	logrus.WithField("666", 1).Infof("hello world")
}
