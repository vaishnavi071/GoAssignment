package initilizers

import (
	"log"
	"os"
)

func GetLogs() {
	logFilePath := "C:/Users/Nagendra/Desktop/GoAssignment/application.log"

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %s", err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("Application finished successfully")
}
