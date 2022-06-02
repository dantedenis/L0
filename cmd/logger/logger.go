package logger

import "log"

type Logger struct {
	ErrorLog *log.Logger `yaml:"-"`
	InfoLog  *log.Logger `yaml:"-"`
}
