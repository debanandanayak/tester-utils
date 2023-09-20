package logger

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

// Logger is a wrapper around log.Logger with the following features:
//   - Supports a prefix
//   - Adds colors to the output
//   - Debug mode (all logs, debug and above)
//   - Quiet mode (only critical logs)
type Logger struct {
	// IsDebug is used to determine whether to emit debug logs.
	IsDebug bool

	// IsQuiet is used to determine whether to emit non-critical logs.
	IsQuiet bool

	logger log.Logger
}

// GetLogger Returns a logger.
func GetLogger(isDebug bool, prefix string) *Logger {
	color.NoColor = false

	prefix = yellowColorize(prefix)
	return &Logger{
		logger:  *log.New(os.Stdout, prefix, 0),
		IsDebug: isDebug,
	}
}

// GetQuietLogger Returns a logger that only emits critical logs. Useful for anti-cheat stages.
func GetQuietLogger(prefix string) *Logger {
	color.NoColor = false

	prefix = yellowColorize(prefix)
	return &Logger{
		logger:  *log.New(os.Stdout, prefix, 0),
		IsDebug: false,
		IsQuiet: true,
	}
}

func (l *Logger) Successf(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}
	msg := successColorize(fstring, args...)
	l.Successln(msg)
}

func (l *Logger) Successln(msg string) {
	if l.IsQuiet {
		return
	}
	msg = successColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Infof(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}
	msg := infoColorize(fstring, args...)
	l.Infoln(msg)
}

func (l *Logger) Infoln(msg string) {
	if l.IsQuiet {
		return
	}
	msg = infoColorize(msg)
	l.logger.Println(msg)
}

// Criticalf is to be used only in anti-cheat stages
func (l *Logger) Criticalf(fstring string, args ...interface{}) {
	if !l.IsQuiet {
		panic("Critical is only for quiet loggers")
	}
	msg := errorColorize(fstring, args...)
	l.Criticalln(msg)
}

// Criticalln is to be used only in anti-cheat stages
func (l *Logger) Criticalln(msg string) {
	if !l.IsQuiet {
		panic("Critical is only for quiet loggers")
	}
	msg = errorColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Errorf(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}
	msg := errorColorize(fstring, args...)
	l.Errorln(msg)
}

func (l *Logger) Errorln(msg string) {
	if l.IsQuiet {
		return
	}
	msg = errorColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Debugf(fstring string, args ...interface{}) {
	if !l.IsDebug {
		return
	}
	msg := debugColorize(fstring, args...)
	l.Debugln(msg)
}

func (l *Logger) Debugln(msg string) {
	if !l.IsDebug {
		return
	}
	msg = debugColorize(msg)
	l.logger.Println(msg)
}

func (l *Logger) Plainln(msg string) {
	l.logger.Println(msg)
}
