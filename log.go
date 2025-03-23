package main

import (
	"log"
	"os"

	"github.com/NatoBoram/BlackCompany/nord"
	"github.com/jwalton/gchalk"
)

type Logger struct {
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	danger *log.Logger
}

var (
	debug  = gchalk.Hex(nord.Nord05.String())
	info   = gchalk.Hex(nord.Nord06.String())
	warn   = gchalk.Hex(nord.Nord13.String())
	danger = gchalk.Hex(nord.Nord11.String())
)

var (
	lflags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
	logger = Logger{
		debug:  log.New(os.Stdout, debug("[DEBUG] "), lflags),
		info:   log.New(os.Stdout, info("[INFO] "), lflags),
		warn:   log.New(os.Stdout, warn("[WARN] "), lflags),
		danger: log.New(os.Stdout, danger("[ERROR] "), lflags),
	}
)

func (l *Logger) Debug(format string, v ...any) {
	l.debug.Printf(format, v...)
}

func (l *Logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.warn.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.danger.Printf(format, v...)
}
