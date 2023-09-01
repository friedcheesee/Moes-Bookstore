package moelog
import (
	"os"
	"log"
)

// opens the log file
func Initiatelog() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) //create a log file, if already exists, appends to the file.
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
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
}