package handlers

import (
	"encoding/json"
	"github.com/TinyMarcus/avito-tech-task/internal/handlers/dtos"
	"github.com/TinyMarcus/avito-tech-task/internal/models"
	"github.com/TinyMarcus/avito-tech-task/internal/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UsersHandler struct {
	repository *repositories.PostgresUserRepository
}

func NewUsersHandler(r *repositories.PostgresUserRepository) *UsersHandler {
	return &UsersHandler{
		repository: r,
	}
}

//go:generate mockgen -source=user_repository.go -destination ./mocks/user_repository.go
type UserRepository interface {
	GetAllUsers() ([]*models.User, error)
	GetUserById(userId string) (*models.User, error)
	CreateUser(name string) error
	DeleteUser(userId string) (*models.User, error)
	AddSegmentToUser(userId, slug, ttl string) error
	TakeSegmentFromUser(userId, slug string) error
	GetActiveSegmentsOfUser(userId string) ([]*dtos.UsersActiveSegments, error)
}

// GetUsersHandler godoc
//
//	@Summary		Получить пользователей
//	@Description	Получить всех пользователей из БД
//	@ID				get-all-users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	    {array} 	dtos.UserDto			"Все пользователи успешно получены"
//	@Failure		500	    {object}	dtos.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users [get]
func (h *UsersHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.repository.GetAllUsers()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dtos.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех пользователей",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	var usersDtos []*dtos.UserDto

	// TODO: некрасиво, стоит переделать
	if users == nil {
		usersDtos = []*dtos.UserDto{}
	}

	for _, val := range users {
		usersDtos = append(usersDtos, dtos.ConvertUserToUserDto(val))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersDtos)
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
//	@Success		200	    {object} 	dtos.UserDto			"Пользователь с данным идентификатором успешно получен"
//	@Failure		400		{object}	dtos.ErrorDto			"Некорректные входные данные"
//	@Failure		404		{object}	dtos.ErrorDto			"Пользователь с данным идентификатором не найден"
//	@Failure		500	    {object}	dtos.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users/{userId} [get]
func (h *UsersHandler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	user, err := h.repository.GetUserById(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		switch err {
		case repositories.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &dtos.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dtos.ErrorDto{
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
//	 	@Param			User	body		dtos.CreateUserDto	    true	"Информация о добавляемом пользователе"
//		@Success		201		{object}	dtos.CreateUserResponseDto						"Пользователь успешно создан"
//		@Failure		400		{object}	dtos.ErrorDto										"Некорректные входные данные"
//		@Failure		500	    {object}	dtos.ErrorDto										"Возникла внутренняя ошибка сервера"
//		@Router			/api/v1/users [post]
func (h *UsersHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user dtos.CreateUserDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dtos.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	id, err := h.repository.CreateUser(user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dtos.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании пользователя",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	createUserResponseDto := dtos.CreateUserResponseDto{
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
//	 	@Param			Информация о добавляемых и удаляемых сегментах	body	dtos.ChangeUserSegmentsDto	    true	"Информация о добавляемых и удаляемых сегментах"
//		@Success		200											"Сегменты пользователя успешно изменены"
//		@Failure		400		{object}	dtos.ErrorDto			"Некорректные входные данные"
//		@Failure		404		{object}	dtos.ErrorDto			"Пользователь с данным идентификатором не найден"
//		@Failure		500	    {object}	dtos.ErrorDto			"Внутренняя ошибка сервера"
//		@Router			/api/v1/users/{userId}/changeSegmentsOfUser [post]
func (h *UsersHandler) ChangeSegmentsOfUserHandler(w http.ResponseWriter, r *http.Request) {
	var userSegment dtos.ChangeUserSegmentsDto
	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userSegment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dtos.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	for _, val := range userSegment.AddToUser {
		err = h.repository.AddSegmentToUser(userId, val.Slug, val.DeadlineDate)
		if err != nil {
			switch err {
			case repositories.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
				errorDto := &dtos.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				json.NewEncoder(w).Encode(errorDto)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &dtos.ErrorDto{
					Error: "Возникла внутренняя ошибка при добавлении новых сегментов пользователю",
				}
				json.NewEncoder(w).Encode(errorDto)
			}
			return
		}
	}

	for _, val := range userSegment.TakeFromUser {
		err = h.repository.TakeSegmentFromUser(userId, val)
		if err != nil {
			switch err {
			case repositories.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
				errorDto := &dtos.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				json.NewEncoder(w).Encode(errorDto)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &dtos.ErrorDto{
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
//	@Success		200		{object}	dtos.UsersActiveSegments		"Активные сегменты пользователя успешно получены"
//	@Failure		404		{object}	dtos.ErrorDto					"Пользователь с данным идентификатором не найден"
//	@Failure		500	    {object}	dtos.ErrorDto					"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/users/{userId}/active [get]
func (h *UsersHandler) GetActiveSegmentsOfUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	usersActiveSegments, err := h.repository.GetActiveSegmentsOfUser(userId)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		switch err {
		case repositories.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &dtos.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dtos.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе активных сегментов пользователя",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersActiveSegments)
}
