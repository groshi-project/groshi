package logger

import (
	"log"
	"os"
)

const baseFlags = log.Ldate | log.Ltime | log.Lmsgprefix

// Info is a logger used to log useful information.
var Info = log.New(os.Stdout, "[info]: ", baseFlags)

// Warning is a logger used to log warnings.
var Warning = log.New(os.Stdout, "[warn]: ", baseFlags)

// Fatal is a logger used to log fatal errors.
var Fatal = log.New(os.Stderr, "[fatal]: ", baseFlags|log.Llongfile)
