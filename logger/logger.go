package logger

import (
	"fmt"
	"log"
)

func LogInfo(message string, err error, tags ...string) {
	fmt.Printf("LogInfo:: %s %v %v \n", message, err, tags)
}

func LogError(message string, err error, tags ...string) {
	fmt.Printf("LogError:: %s %v %v \n", message, err, tags)
}

func LogFatal(message string, err error, tags ...string) {
	formattedMessage := fmt.Sprintf("LogFatal:: %s %v\n", message, tags)
	log.Fatal(formattedMessage, err)
}
