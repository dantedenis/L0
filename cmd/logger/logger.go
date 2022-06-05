package logger

import (
	"log"
	"os"
)

type Logger struct {
	ErrorLog *log.Logger `yaml:"-"`
	InfoLog  *log.Logger `yaml:"-"`
}

func NewLogger() *Logger {
	return &Logger{
		ErrorLog: log.New(os.Stderr, "\033[31mERROR\t\033[0m", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLog:  log.New(os.Stdout, "\033[32mINFO\t\033[0m", log.Ldate|log.Ltime),
	}
}
