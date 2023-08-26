package main

import (
	"dynamic-user-segmentation-service/internal/handlers"
	"dynamic-user-segmentation-service/internal/utils"
	"net/http"
	"os"
)

const appName = "dynamic-user-segmentation-service"

func main() {
	port := os.Getenv("PORT")
	logger := utils.CreateLogger()
	defer logger.Sync()

	r := handlers.Router(logger)

	logger.Info("Server is started on port ", port)
	logger.Fatal(http.ListenAndServe(":"+port, r))
}
