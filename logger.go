package tester_utils

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func colorize(colorToUse color.Attribute, fstring string, args ...interface{}) string {
	return color.New(colorToUse).SprintfFunc()(fstring, args...)
}

func debugColorize(fstring string, args ...interface{}) string {
	return colorize(color.FgCyan, fstring, args...)
}

func infoColorize(fstring string, args ...interface{}) string {
	return colorize(color.FgHiBlue, fstring, args...)
}

func successColorize(fstring string, args ...interface{}) string {
	return colorize(color.FgHiGreen, fstring, args...)
}

func errorColorize(fstring string, args ...interface{}) string {
	return colorize(color.FgHiRed, fstring, args...)
}

func yellowColorize(fstring string, args ...interface{}) string {
	return colorize(color.FgYellow, fstring, args...)
}

type Logger struct {
	logger  log.Logger
	isDebug bool
	isQuiet bool // Only CRITICAL logs
}

func getLogger(isDebug bool, prefix string) *Logger {
	color.NoColor = false

	prefix = yellowColorize(prefix)
	return &Logger{
		logger:  *log.New(os.Stdout, prefix, 0),
		isDebug: isDebug,
	}
}

func getQuietLogger(prefix string) *Logger {
	color.NoColor = false

	prefix = yellowColorize(prefix)
	return &Logger{
		logger:  *log.New(os.Stdout, prefix, 0),
		isDebug: false,
		isQuiet: true,
	}
}

func (l *Logger) Successf(fstring string, args ...interface{}) {
	if l.isQuiet {
		return
	}
	msg := successColorize(fstring, args...)
	l.Successln(msg)
}

func (l *Logger) Successln(msg string) {
	if l.isQuiet {
		return
	}
	msg = successColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Infof(fstring string, args ...interface{}) {
	if l.isQuiet {
		return
	}
	msg := infoColorize(fstring, args...)
	l.Infoln(msg)
}

func (l *Logger) Infoln(msg string) {
	if l.isQuiet {
		return
	}
	msg = infoColorize(msg)
	l.logger.Println(msg)
}

// Criticalf is to be used only in anti-cheat stages
func (l *Logger) Criticalf(fstring string, args ...interface{}) {
	if !l.isQuiet {
		panic("Critical is only for quiet loggers")
	}
	msg := errorColorize(fstring, args...)
	l.Criticalln(msg)
}

// Criticalln is to be used only in anti-cheat stages
func (l *Logger) Criticalln(msg string) {
	if !l.isQuiet {
		panic("Critical is only for quiet loggers")
	}
	msg = errorColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Errorf(fstring string, args ...interface{}) {
	if l.isQuiet {
		return
	}
	msg := errorColorize(fstring, args...)
	l.Errorln(msg)
}

func (l *Logger) Errorln(msg string) {
	if l.isQuiet {
		return
	}
	msg = errorColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Debugf(fstring string, args ...interface{}) {
	if !l.isDebug {
		return
	}
	msg := debugColorize(fstring, args...)
	l.Debugln(msg)
}

func (l *Logger) Debugln(msg string) {
	if !l.isDebug {
		return
	}
	msg = debugColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Plainln(msg string) {
	l.logger.Println(msg)
}
