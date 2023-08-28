package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/TinyMarcus/avito-tech-task/internal/handlers/dto"
	"github.com/TinyMarcus/avito-tech-task/internal/models"
	"github.com/TinyMarcus/avito-tech-task/internal/repositories"
)

type UsersHandler struct {
	repository UserRepository
}

func NewUsersHandler(r UserRepository) *UsersHandler {
	return &UsersHandler{
		repository: r,
	}
}

//go:generate mockgen -source=user_repository.go -destination ./mocks/user_repository.go
type UserRepository interface {
	GetAllUsers() ([]*models.User, error)
	GetUserById(userId int) (*models.User, error)
	CreateUser(name string) (int, error)
	AddSegmentToUser(userId int, slug, ttl string) error
	TakeSegmentFromUser(userId int, slug string) error
	GetActiveSegmentsOfUser(userId int) (*dto.UsersActiveSegments, error)
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
		errorDto := &dto.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех пользователей",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var usersDtos []*dto.UserDto

	// TODO: некрасиво, стоит переделать
	if users == nil {
		usersDtos = []*dto.UserDto{}
	}

	for _, val := range users {
		usersDtos = append(usersDtos, dto.ConvertUserToUserDto(val))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(usersDtos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
			errorDto := &dto.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dto.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе пользователя",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	var user dto.CreateUserDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные входные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	id, err := h.repository.CreateUser(user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании пользователя",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	createUserResponseDto := dto.CreateUserResponseDto{
		Id: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createUserResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	var userSegment dto.ChangeUserSegmentsDto
	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userSegment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные входные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	for _, val := range userSegment.AddToUser {
		err = h.repository.AddSegmentToUser(userId, val.Slug, val.DeadlineDate)
		if err != nil {
			switch err {
			case repositories.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
				errorDto := &dto.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				err = json.NewEncoder(w).Encode(errorDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &dto.ErrorDto{
					Error: "Возникла внутренняя ошибка при добавлении новых сегментов пользователю",
				}
				err = json.NewEncoder(w).Encode(errorDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
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
				errorDto := &dto.ErrorDto{
					Error: "Пользователь с таким идентификатором не найден",
				}
				err = json.NewEncoder(w).Encode(errorDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			default:
				w.WriteHeader(http.StatusInternalServerError)
				errorDto := &dto.ErrorDto{
					Error: "Возникла внутренняя ошибка при удалении сегментов пользователя",
				}
				err = json.NewEncoder(w).Encode(errorDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
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
			errorDto := &dto.ErrorDto{
				Error: "Пользователь с таким идентификатором не найден",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dto.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе активных сегментов пользователя",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(usersActiveSegments)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
