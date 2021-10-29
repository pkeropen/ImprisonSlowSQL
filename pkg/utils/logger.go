package utils

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

type LogFormatter struct{}

// Format Log输出格式
func (m *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s][%s][%s:%d] %s\n", timestamp, strings.ToUpper(entry.Level.String()), fName, entry.Caller.Line, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s][%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
