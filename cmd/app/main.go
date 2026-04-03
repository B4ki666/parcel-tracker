package main

import (
	"log"
	"os"
	"parcel_tracker/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	server := server.NewServer(logger)

	err := server.Start()
	if err != nil {
		logger.Fatal(err)
	}

}
