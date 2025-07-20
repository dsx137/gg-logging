package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type GeneralFormatter struct{}

func (f *GeneralFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &strings.Builder{}
	if entry.Buffer != nil {
		b.Write(entry.Buffer.Bytes())
	}

	now := time.Now().Format("15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	showPath, line := "unknown_file", 0
	if entry.Caller != nil {
		showPath = entry.Caller.File
		line = entry.Caller.Line
	}

	if logrus.GetLevel() > logrus.InfoLevel {
		fmt.Fprintf(b, "[%s %5s] [%s:%d]: %s\n", now, level, showPath, line, message)
	} else {
		fmt.Fprintf(b, "[%s %5s]: %s\n", now, level, message)
	}
	return []byte(b.String()), nil
}

func Init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&GeneralFormatter{})
}
