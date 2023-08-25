package handlers

import (
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/models"
	"dynamic-user-segmentation-service/internal/repositories"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	users, err := userRepository.GetAllUsers()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("failed to get users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех пользователей",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	json.NewEncoder(w).Encode(users)
	w.WriteHeader(http.StatusOK)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	user, err := userRepository.GetUserById(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("failed to get user: %v", err)
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе пользователя",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	var user models.CreateUserDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	err = userRepository.CreateUser(user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании пользователя",
		}
		json.NewEncoder(w).Encode(errorDto)
		log.Printf("failed to create user: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ChangeSegmentsOfUserHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	var userSegment models.ChangeUserSegmentsDto
	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userSegment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	for _, val := range userSegment.AddToUser {
		err = userRepository.AddSegmentToUser(userId, val.Slug, val.DeadlineDate)
		if err != nil {
			log.Printf("failed to add segment to user: %v\n", err)
			switch err {
			case errors.RecordNotFound:
				w.WriteHeader(http.StatusNotFound)
				errorDto := &models.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				json.NewEncoder(w).Encode(errorDto)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &models.ErrorDto{
					Error: "Возникла внутренняя ошибка при добавлении новых сегментов пользователю",
				}
				json.NewEncoder(w).Encode(errorDto)
			}
			return
		}
	}

	for _, val := range userSegment.TakeFromUser {
		err = userRepository.TakeSegmentFromUser(userId, val)
		if err != nil {
			log.Printf("failed to take segment from user: %v\n", err)
			switch err {
			case errors.RecordNotFound:
				w.WriteHeader(http.StatusNotFound)
				errorDto := &models.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				json.NewEncoder(w).Encode(errorDto)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &models.ErrorDto{
					Error: "Возникла внутренняя ошибка при удалении сегментов пользователя",
				}
				json.NewEncoder(w).Encode(errorDto)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func GetActiveSegmentsOfUser(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	usersActiveSegments, err := userRepository.GetActiveSegmentsOfUser(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("failed to get user's active segments: %v", err)
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе активных сегментов пользователя",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	json.NewEncoder(w).Encode(usersActiveSegments)
	w.WriteHeader(http.StatusOK)
}
