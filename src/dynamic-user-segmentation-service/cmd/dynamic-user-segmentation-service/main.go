package main

import (
	"dynamic-user-segmentation-service/internal/handlers"
	"dynamic-user-segmentation-service/internal/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	utils.InitConfiguration()
	r := handlers.Router()

	log.Println("server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
