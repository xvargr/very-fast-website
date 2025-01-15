package logger

import (
	"fmt"
	"time"
)

type Severity int

const (
	SeverityNormal Severity = iota
	SeverityDebug
	SeverityInfo
	SeverityWarning
	SeverityError
)

var severityStrings = map[Severity]string{
	SeverityNormal:  "NORMAL",
	SeverityDebug:   "DEBUG",
	SeverityInfo:    "INFO",
	SeverityWarning: "WARNING",
	SeverityError:   "ERROR",
}

var colors = map[Severity]string{
	SeverityNormal:  "\033[0;37m",
	SeverityDebug:   "\033[0;35m",
	SeverityInfo:    "\033[0;34m",
	SeverityWarning: "\033[0;33m",
	SeverityError:   "\033[0;31m",
}

func Console(severity Severity, message string) {
	fmt.Println(
		colors[severity]+
			time.Now().Format(time.DateTime),
		message,
		"\033[0m",
	)
}
