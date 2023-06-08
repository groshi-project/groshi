package loggers

import (
	"log"
	"os"
)

const baseFlags = log.Ldate | log.Ltime | log.Lmsgprefix

var Info = log.New(os.Stdout, "[info]: ", baseFlags)
var Error = log.New(os.Stderr, "[error]: ", baseFlags|log.Llongfile)
