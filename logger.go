package main

import (
	"fmt"
	"io"
)

type logLevel uint8

const (
  DEBUG logLevel = iota
  INFO
  WARN
  ERROR
)

type Logger struct {
  output io.Writer
  level logLevel
}

func newLogger(out io.Writer, level logLevel) *Logger {
  return &Logger{out, level}
}

func (l *Logger) log(level logLevel, msg string, args ...interface{}) {
	if(level >= l.level) {
		fmt.Fprintf(l.output, msg  + "\n", args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)	
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)	
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)	
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)	
}