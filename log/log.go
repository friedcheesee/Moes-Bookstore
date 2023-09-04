package moelog

import (
	"fmt"
	"log"
	"os"
)

// opens the log file
func Initiatelog() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) //create a log file, if already exists, appends to the file.
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
	}
	fmt.Println("Check app.log in the root directory for detailed logs/previously stored logs")
	log.SetOutput(logFile)
	return logFile
}

// logs errors
func CheckError(err error) { // to log errors where ever necessary
	if err != nil {
		log.Println("Error:", err)
		//panic(err)
	}
}

// to log events
func LogEvent(message string) {
	// Log the message
	log.Println(message)
	fmt.Println(message)
}
