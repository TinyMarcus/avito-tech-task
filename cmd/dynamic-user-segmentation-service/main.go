package main

import (
	"github.com/TinyMarcus/avito-tech-task/internal/config"
	"github.com/TinyMarcus/avito-tech-task/internal/db"
	"github.com/TinyMarcus/avito-tech-task/internal/handlers"
	"github.com/TinyMarcus/avito-tech-task/internal/logger"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"net/http"
)

// @title       Dynamic User Segmentation Service
// @version     1.0
// @description Dynamic User Segmentation Service

func main() {
	config, err := config.New()
	logger := logger.CreateLogger(config.Log)
	defer logger.Sync()
	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	db, err := db.CreateConnection(config.Db)
	defer db.Close()
	if err != nil {
		logger.Fatalf("Error while connecting to database: %v", err)
	}

	port := config.Port
	r := handlers.Router(logger, db)

	logger.Info("Server is started on port ", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatalf("Error while starting server: %v", err)
	}
}
