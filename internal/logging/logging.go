package logging

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type GeneralFormatter struct {
	ShowFields bool
}

func (f *GeneralFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &strings.Builder{}
	if entry.Buffer != nil {
		b.Write(entry.Buffer.Bytes())
	}

	now := entry.Time.Format("15:04:05")
	level := strings.ToUpper(entry.Level.String())

	showPath, line := "unknown_file", 0
	if entry.Caller != nil {
		showPath = entry.Caller.File
		line = entry.Caller.Line
	}

	msg := entry.Message
	if f.ShowFields {
		if fieldText := formatFields(entry.Data); fieldText != "" {
			msg += " " + fieldText
		}
	}

	if logrus.GetLevel() > logrus.InfoLevel {
		fmt.Fprintf(b, "[%s %5s] [%s:%d]: %s\n", now, level, showPath, line, msg)
	} else {
		fmt.Fprintf(b, "[%s %5s]: %s\n", now, level, msg)
	}
	return []byte(b.String()), nil
}

func formatFields(fields logrus.Fields) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, formatValue(fields[k])))
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

func formatValue(v any) string {
	s := fmt.Sprint(v)
	if strings.ContainsAny(s, " \t\n\"=") {
		return strconv.Quote(s)
	}
	return s
}

func Init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&GeneralFormatter{
		ShowFields: true,
	})
}
