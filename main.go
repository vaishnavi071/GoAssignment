package main

import (
	"GoAssignment/internal/database"
	initilizers "GoAssignment/internal/initializers"
	"GoAssignment/internal/student"

	transportHTTP "GoAssignment/internal/transport"

	log "github.com/sirupsen/logrus"
)

func init() {
	initilizers.LoadEnvVariables()
}

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting Up Our APP")

	var err error

	db, err := database.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}

	stdentService := student.NewService(db)
	handler := transportHTTP.NewHandler(stdentService)
	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
