package main

import (
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/TinyMarcus/avito-tech-task/internal/config"
	"github.com/TinyMarcus/avito-tech-task/internal/db"
	"github.com/TinyMarcus/avito-tech-task/internal/handlers"
	"github.com/TinyMarcus/avito-tech-task/internal/logger"
	"github.com/TinyMarcus/avito-tech-task/internal/repositories"
)

// @title       Dynamic User Segmentation Service
// @version     1.0
// @description Dynamic User Segmentation Service

func main() {
	config, err := config.New()
	logger := logger.CreateLogger(config.Log)

	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Errorf("Error while syncing logger: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	db, err := db.CreateConnection(config.Db)

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Errorf("Error while closing connection to db: %v", err)
		}
	}()

	if err != nil {
		logger.Fatalf("Error while connecting to database: %v", err)
	}

	hr := repositories.NewHistoryRepository(db)
	ur := repositories.NewUserRepository(db, hr)
	sr := repositories.NewSegmentRepository(db)
	r := handlers.Router(logger, ur, sr)

	port := config.Port
	logger.Info("Server is started on port ", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatalf("Error while starting server: %v", err)
	}
}
