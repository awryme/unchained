package clilog

import (
	"log"
	"os"
)

var Output = os.Stdout

var logger = log.New(Output, "", 0)

func Log(v ...any) {
	logger.Println(v...)
}

func Logf(format string, v ...any) {
	logger.Printf(format, v...)
}
