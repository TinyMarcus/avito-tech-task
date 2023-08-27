package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/TinyMarcus/avito-tech-task/internal/handlers/dtos"
	"github.com/TinyMarcus/avito-tech-task/internal/repositories"
)

type SegmentsHandler struct {
	repository *repositories.PostgresSegmentRepository
}

func NewSegmentsHandler(r *repositories.PostgresSegmentRepository) *SegmentsHandler {
	return &SegmentsHandler{
		repository: r,
	}
}

//go:generate mockgen -source=segment_repository.go -destination ./mocks/segment_repository.go
type SegmentRepository interface {
	GetAllSegments() ([]*dtos.SegmentDto, error)
	GetSegmentBySlug(slug string) (*dtos.SegmentDto, error)
	CreateSegment(slug, description string) (string, error)
	UpdateSegment(slug, description string) (*dtos.SegmentDto, error)
	DeleteSegment(slug string) (*dtos.SegmentDto, error)
}

// GetSegmentsHandler godoc
//
//	@Summary		Получить семгенты
//	@Description	Получить все сегменты из БД
//	@ID				get-all-segments
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Success		200	    {array} 	dtos.SegmentDto			"Все сегменты успешно получены"
//	@Failure		500	    {object}	dtos.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments [get]
func (h *SegmentsHandler) GetSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	segments, err := h.repository.GetAllSegments()

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dtos.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех сегментов",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var segmentsDtos []*dtos.SegmentDto

	// TODO: некрасиво, стоит переделать
	if segments == nil {
		segmentsDtos = []*dtos.SegmentDto{}
	}

	for _, val := range segments {
		segmentsDtos = append(segmentsDtos, dtos.ConvertSegmentToSegmentDto(val))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(segmentsDtos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetSegmentBySlugHandler godoc
//
//	@Summary		Получить сегмент
//	@Description	Получить сегмент из БД по названию
//	@ID				get-segment-by-name
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string					true	"Название сегмента"
//	@Success		200	    {object} 	dtos.SegmentDto			"Сегмент с данным названием успешно получен"
//	@Failure		400		{object}	dtos.ErrorDto			"Некорректные входные данные"
//	@Failure		404		{object}	dtos.ErrorDto			"Сегмент с данным названием не найден"
//	@Failure		500	    {object}	dtos.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments/{slug} [get]
func (h *SegmentsHandler) GetSegmentBySlugHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	segment, err := h.repository.GetSegmentBySlug(slug)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		switch err {
		case repositories.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &dtos.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dtos.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе сегмента по названию",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ConvertSegmentToSegmentDto(segment))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CreateSegmentHandler godoc
//
//		@Summary		Добавить сегмент
//		@Description	Добавить сегмент в БД
//		@ID				create-segment
//		@Tags			segments
//		@Accept			json
//		@Produce		json
//	 	@Param			Segment 	body	dtos.CreateOrUpdateSegmentDto	    true	"Информация о добавляемом сегменте"
//		@Success		201		{object}	dtos.CreateSegmentResponseDto		"Сегмент успешно создан"
//		@Failure		400		{object}	dtos.ErrorDto						"Некорректные входные данные"
//		@Failure		500	    {object}	dtos.ErrorDto						"Возникла внутренняя ошибка сервера"
//		@Router			/api/v1/segments [post]
func (h *SegmentsHandler) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	var segment dtos.CreateOrUpdateSegmentDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dtos.ErrorDto{
			Error: "Некорректные входные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	slug, err := h.repository.CreateSegment(segment.Slug, segment.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dtos.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании сегмента",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	createResponseDto := dtos.CreateSegmentResponseDto{
		Slug: slug,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// UpdateSegmentHandler godoc
//
//		@Summary		Обновить сегмент
//		@Description	Обновить сегмент в БД
//		@ID				update-segment
//		@Tags			segments
//		@Accept			json
//		@Produce		json
//		@Param			slug	path		string					true		"Название сегмента"
//	 	@Param			Информация о сегменте	body	dtos.CreateOrUpdateSegmentDto	    true	"Информация о добавляемом сегменте"
//		@Success		200		{object}	dtos.UpdateSegmentResponseDto		"Сегмент с данным названием успешно обновлен"
//		@Failure		400		{object}	dtos.ErrorDto						"Некорректные входные данные"
//		@Failure		404		{object}	dtos.ErrorDto						"Сегмент с данным названием не найден"
//		@Failure		500	    {object}	dtos.ErrorDto						"Возникла внутренняя ошибка сервреа"
//		@Router			/api/v1/segments/{slug} [put]
func (h *SegmentsHandler) UpdateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	var segment dtos.CreateOrUpdateSegmentDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dtos.ErrorDto{
			Error: "Некорректные входные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	updated, err := h.repository.UpdateSegment(slug, segment.Description)
	if err != nil {
		switch err {
		case repositories.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &dtos.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dtos.ErrorDto{
				Error: "Возникла внутренняя ошибка при обновлении сегмента",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		return
	}

	updateResponseDto := dtos.UpdateSegmentResponseDto{
		Id:          updated.Id,
		Slug:        updated.Slug,
		Description: updated.Description,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updateResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// DeleteSegmentHandler godoc
//
//	@Summary		Удалить сегмент
//	@Description	Удалить сегмент в БД
//	@ID				delete-segment
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string					true		"Название сегмента"
//	@Success		204														"Сегмент с данным названием успешно удален"
//	@Failure		400		{object}	dtos.ErrorDto						"Некорректные входные данные"
//	@Failure		404		{object}	dtos.ErrorDto						"Сегмент с данным названием не найден"
//	@Failure		500	    {object}	dtos.ErrorDto						"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments/{slug} [delete]
func (h *SegmentsHandler) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	_, err := h.repository.DeleteSegment(slug)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		switch err {
		case repositories.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &dtos.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &dtos.ErrorDto{
				Error: "Возникла внутренняя ошибка при удалении сегмента",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
