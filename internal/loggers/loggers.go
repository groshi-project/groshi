package loggers

import (
	"log"
	"os"
)

const flags = log.Ldate | log.Ltime | log.Lmsgprefix

var Info = log.New(os.Stdout, "[info]: ", flags)
var Error = log.New(os.Stderr, "[error]: ", flags)
var Fatal = log.New(os.Stderr, "[fatal]: ", flags)
