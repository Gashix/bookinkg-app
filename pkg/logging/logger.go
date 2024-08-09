package logging

import (
	"fmt"
	"log"
)

var logger = log.Default()

func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
