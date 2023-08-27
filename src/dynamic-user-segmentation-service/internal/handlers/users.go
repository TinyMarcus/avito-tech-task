package handlers

import (
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/models"
	"dynamic-user-segmentation-service/internal/repositories"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetUsersHandler godoc
//
//	@Summary		Получить пользователей
//	@Description	Получить всех пользователей из БД
//	@ID				get-all-users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	    {array} 	models.User				"Все пользователи успешно получены"
//	@Failure		500	    {object}	models.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users [get]
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	users, err := userRepository.GetAllUsers()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех пользователей",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	// TODO: некрасиво, стоит переделать
	if users == nil {
		users = []*models.User{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUserByIdHandler godoc
//
//	@Summary		Получить пользователя
//	@Description	Получить пользователя из БД по идентификатору
//	@ID				get-user-by-id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int						true	"Идентификатор пользователя"
//	@Success		200	    {object} 	models.User				"Пользователь с данным идентификатором успешно получен"
//	@Failure		400		{object}	models.ErrorDto			"Некорректные входные данные"
//	@Failure		404		{object}	models.ErrorDto			"Пользователь с данным идентификатором не найден"
//	@Failure		500	    {object}	models.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users/{userId} [get]
func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	user, err := userRepository.GetUserById(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// CreateUserHandler godoc
//
//		@Summary		Добавить пользователя
//		@Description	Добавить пользователя в БД
//		@ID				create-user
//		@Tags			users
//		@Accept			json
//		@Produce		json
//	 	@Param			Информация о пользователе	body	models.CreateUserDto	    true	"Информация о добавляемом пользователе"
//		@Success		201		{object}	models.CreateUserResponseDto						"Пользователь успешно создан"
//		@Failure		400		{object}	models.ErrorDto										"Некорректные входные данные"
//		@Failure		500	    {object}	models.ErrorDto										"Возникла внутренняя ошибка сервера"
//		@Router			/api/v1/users [post]
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

	id, err := userRepository.CreateUser(user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании пользователя",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	createUserResponseDto := models.CreateUserResponseDto{
		Id: id,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createUserResponseDto)
}

// ChangeSegmentsOfUserHandler godoc
//
//		@Summary		Изменить сегменты пользователя
//		@Description	Добавить и удалить у пользователя указанные сегменты
//		@ID				change-segments-of-user
//		@Tags			users
//		@Accept			json
//		@Produce		json
//		@Param			userId	path		int						true	"Идентификатор пользователя"
//	 	@Param			Информация о добавляемых и удаляемых сегментах	body	models.ChangeUserSegmentsDto	    true	"Информация о добавляемых и удаляемых сегментах"
//		@Success		200											"Сегменты пользователя успешно изменены"
//		@Failure		400		{object}	models.ErrorDto			"Некорректные входные данные"
//		@Failure		404		{object}	models.ErrorDto			"Пользователь с данным идентификатором не найден"
//		@Failure		500	    {object}	models.ErrorDto			"Внутренняя ошибка сервера"
//		@Router			/api/v1/users/{userId}/changeSegmentsOfUser [post]
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

// GetActiveSegmentsOfUser godoc
//
//	@Summary		Получить активные сегменты пользователя
//	@Description	Получить из БД активные сегменты пользователя, в которых он участвует на момент запроса
//	@ID				get-active-segments-of-user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int					true		"Идентификатор пользователя"
//	@Success		200		{object}	models.UsersActiveSegments		"Активные сегменты пользователя успешно получены"
//	@Failure		404		{object}	models.ErrorDto					"Пользователь с данным идентификатором не найден"
//	@Failure		500	    {object}	models.ErrorDto					"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users/{userId}/active [get]
func GetActiveSegmentsOfUser(w http.ResponseWriter, r *http.Request) {
	userRepository := repositories.PostgresUserRepository{}

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	usersActiveSegments, err := userRepository.GetActiveSegmentsOfUser(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersActiveSegments)
}
