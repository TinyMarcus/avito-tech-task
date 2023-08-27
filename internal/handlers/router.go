package handlers

import (
	_ "github.com/TinyMarcus/avito-tech-task/api"
	"github.com/TinyMarcus/avito-tech-task/internal/handlers/middlewares"
	"github.com/TinyMarcus/avito-tech-task/internal/repositories"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func Router(logger *zap.SugaredLogger, db *sqlx.DB) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.Use(middlewares.LoggerMiddleware(logger))

	segmentsHandler := NewSegmentsHandler(repositories.NewSegmentRepository(db))
	router.HandleFunc("/api/v1/segments", segmentsHandler.GetSegmentsHandler).Methods("GET")
	router.HandleFunc("/api/v1/segments/{slug}", segmentsHandler.GetSegmentBySlugHandler).Methods("GET")
	router.HandleFunc("/api/v1/segments", segmentsHandler.CreateSegmentHandler).Methods("POST")
	router.HandleFunc("/api/v1/segments/{slug}", segmentsHandler.UpdateSegmentHandler).Methods("PUT")
	router.HandleFunc("/api/v1/segments/{slug}", segmentsHandler.DeleteSegmentHandler).Methods("DELETE")

	hr := repositories.NewHistoryRepository(db)
	usersHandler := NewUsersHandler(repositories.NewUserRepository(db, hr))
	router.HandleFunc("/api/v1/users", usersHandler.GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/v1/users/{userId}", usersHandler.GetUserByIdHandler).Methods("GET")
	router.HandleFunc("/api/v1/users", usersHandler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/changeSegmentsOfUser", usersHandler.ChangeSegmentsOfUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{userId}/active", usersHandler.GetActiveSegmentsOfUser).Methods("GET")

	return router
}
