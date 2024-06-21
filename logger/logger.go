package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

func colorize(colorToUse color.Attribute, fstring string, args ...interface{}) []string {
	msg := fmt.Sprintf(fstring, args...)
	lines := strings.Split(msg, "\n")
	colorizedLines := make([]string, len(lines))

	for i, line := range lines {
		colorizedLines[i] = color.New(colorToUse).SprintFunc()(line)
	}

	return colorizedLines
}

func debugColorize(fstring string, args ...interface{}) []string {
	return colorize(color.FgCyan, fstring, args...)
}

func infoColorize(fstring string, args ...interface{}) []string {
	return colorize(color.FgHiBlue, fstring, args...)
}

func successColorize(fstring string, args ...interface{}) []string {
	return colorize(color.FgHiGreen, fstring, args...)
}

func errorColorize(fstring string, args ...interface{}) []string {
	return colorize(color.FgHiRed, fstring, args...)
}

func yellowColorize(fstring string, args ...interface{}) []string {
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

	// Prefix is the prefix to be used for all logs.
	Prefix string

	logger log.Logger
}

// GetLogger Returns a logger.
func GetLogger(isDebug bool, prefix string) *Logger {
	color.NoColor = false

	coloredPrefix := yellowColorize(prefix)[0]
	return &Logger{
		logger:  *log.New(os.Stdout, coloredPrefix, 0),
		IsDebug: isDebug,
		Prefix:  prefix,
	}
}

// GetLoggerWithAppendedPrefix Returns a logger after appending a prefix to the existing prefix.
func (l *Logger) GetLoggerWithAppendedPrefix(prefix string) *Logger {
	color.NoColor = false

	coloredPrefix := yellowColorize(l.Prefix + fmt.Sprintf("[%s] ", prefix))[0]
	return &Logger{
		logger:  *log.New(os.Stdout, coloredPrefix, 0),
		IsDebug: l.IsDebug,
		Prefix:  prefix,
	}
}

// GetQuietLogger Returns a logger that only emits critical logs. Useful for anti-cheat stages.
func GetQuietLogger(prefix string) *Logger {
	color.NoColor = false

	coloredPrefix := yellowColorize(prefix)[0]
	return &Logger{
		logger:  *log.New(os.Stdout, coloredPrefix, 0),
		IsDebug: false,
		IsQuiet: true,
		Prefix:  prefix,
	}
}

func (l *Logger) Successf(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}

	for _, line := range successColorize(fstring, args...) {
		l.logger.Println(line)
	}
}

func (l *Logger) Successln(msg string) {
	if l.IsQuiet {
		return
	}
	for _, line := range successColorize(msg) {
		l.logger.Println(line)
	}
}

func (l *Logger) Infof(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}

	for _, line := range infoColorize(fstring, args...) {
		l.logger.Println(line)
	}
}

func (l *Logger) Infoln(msg string) {
	if l.IsQuiet {
		return
	}

	for _, line := range infoColorize(msg) {
		l.logger.Println(line)
	}
}

// Criticalf is to be used only in anti-cheat stages
func (l *Logger) Criticalf(fstring string, args ...interface{}) {
	if !l.IsQuiet {
		panic("Critical is only for quiet loggers")
	}

	for _, line := range errorColorize(fstring, args...) {
		l.logger.Println(line)
	}
}

// Criticalln is to be used only in anti-cheat stages
func (l *Logger) Criticalln(msg string) {
	if !l.IsQuiet {
		panic("Critical is only for quiet loggers")
	}

	for _, line := range errorColorize(msg) {
		l.logger.Println(line)
	}
}

func (l *Logger) Errorf(fstring string, args ...interface{}) {
	if l.IsQuiet {
		return
	}

	for _, line := range errorColorize(fstring, args...) {
		l.logger.Println(line)
	}
}

func (l *Logger) Errorln(msg string) {
	if l.IsQuiet {
		return
	}

	for _, line := range errorColorize(msg) {
		l.logger.Println(line)
	}
}

func (l *Logger) Debugf(fstring string, args ...interface{}) {
	if !l.IsDebug {
		return
	}

	for _, line := range debugColorize(fstring, args...) {
		l.logger.Println(line)
	}
}

func (l *Logger) Debugln(msg string) {
	if !l.IsDebug {
		return
	}

	for _, line := range debugColorize(msg) {
		l.logger.Println(line)
	}
}

func (l *Logger) Plainf(fstring string, args ...interface{}) {
	formattedString := fmt.Sprintf(fstring, args...)

	for _, line := range strings.Split(formattedString, "\n") {
		l.logger.Println(line)
	}
}

func (l *Logger) Plainln(msg string) {
	lines := strings.Split(msg, "\n")

	for _, line := range lines {
		l.logger.Println(line)
	}
}
