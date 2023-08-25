package handlers

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/segments", GetSegmentsHandler).Methods("GET")
	router.HandleFunc("/api/v1/segments/{slug}", GetSegmentBySlugHandler).Methods("GET")
	router.HandleFunc("/api/v1/segments", CreateSegmentHandler).Methods("POST")
	router.HandleFunc("/api/v1/segments/{slug}", UpdateSegmentHandler).Methods("PATCH")
	router.HandleFunc("/api/v1/segments/{slug}", DeleteSegmentHandler).Methods("DELETE")

	router.HandleFunc("/api/v1/users", GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/v1/users/{userId}", GetUserByIdHandler).Methods("GET")
	router.HandleFunc("/api/v1/users", CreateUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/changeSegmentsOfUser", ChangeSegmentsOfUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/active", GetActiveSegmentsOfUser).Methods("GET")

	return router
}
