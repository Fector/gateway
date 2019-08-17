package conn

import (
	"fmt"
	"log"
)

type Logger struct {
	Prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{Prefix: prefix}
}

func (l Logger) Println(v ...interface{}) {
	log.Println(fmt.Sprintf("[gateway:%s] ", fmt.Sprintln(v)))
}
