package main

import (
	"log"
	"os"
	"parcel_tracker/internal/db"
	"parcel_tracker/internal/repository"
	"parcel_tracker/internal/server"
	"parcel_tracker/internal/service"
)

func main() {
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//Подключаем БД
	db, err := db.NewSQLiteDB("parcels.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Создаем repository
	parcelRepo := repository.NewParcelRepository(db)
	parcelService := service.NewParcelService(parcelRepo)

	//Передаем зависимости в сервер
	server := server.NewServer(logger, parcelService)

	//Запускаем сервер
	err = server.Start()
	if err != nil {
		logger.Fatal(err)
	}

}
