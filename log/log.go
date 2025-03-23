package log

import (
	"log"
	"os"

	"github.com/NatoBoram/BlackCompany/nord"
	"github.com/jwalton/gchalk"
)

const lflags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix

var (
	debug  = log.New(os.Stdout, colourDebug("[DEBUG] "), lflags)
	info   = log.New(os.Stdout, colourInfo("[INFO] "), lflags)
	warn   = log.New(os.Stdout, colourWarn("[WARN] "), lflags)
	danger = log.New(os.Stdout, colourDanger("[ERROR] "), lflags)
	fatal  = log.New(os.Stdout, colourFatal("[FATAL] "), lflags)
)

var (
	colourDebug  = gchalk.Hex(nord.Nord05.String())
	colourInfo   = gchalk.Hex(nord.Nord06.String())
	colourWarn   = gchalk.Hex(nord.Nord13.String())
	colourDanger = gchalk.Hex(nord.Nord12.String())
	colourFatal  = gchalk.Hex(nord.Nord11.String())
)

func Debug(format string, v ...any) {
	debug.Printf(format, v...)
}

func Info(format string, v ...any) {
	info.Printf(format, v...)
}

func Warn(format string, v ...any) {
	warn.Printf(format, v...)
}

func Error(format string, v ...any) {
	danger.Printf(format, v...)
}

func Fatal(format string, v ...any) {
	fatal.Fatalf(format, v...)
}
