package main

import "log"

type logger struct {
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func NewLogger(cfg config) *logger {
	return &logger{
		info:    log.New(cfg.LogTo, "INFO: ", log.LstdFlags),
		warning: log.New(cfg.LogTo, "WARN: ", log.LstdFlags),
		error:   log.New(cfg.LogTo, "ERROR: ", log.LstdFlags),
	}
}

func (l *logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}
func (l *logger) Warning(format string, v ...any) {
	l.warning.Printf(format, v...)
}
func (l *logger) Error(format string, v ...any) {
	l.error.Printf(format, v...)
}
