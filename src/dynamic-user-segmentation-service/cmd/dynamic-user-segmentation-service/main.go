package main

import (
	"dynamic-user-segmentation-service/internal/handlers"
	"dynamic-user-segmentation-service/internal/utils"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	logger := utils.CreateLogger()
	defer logger.Sync()

	r := handlers.Router(logger)

	logger.Info("Server is started on port ", port)
	logger.Fatal(http.ListenAndServe(":"+port, r))
}
