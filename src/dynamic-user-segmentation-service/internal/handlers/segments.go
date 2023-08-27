package handlers

import (
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/repositories"
	"encoding/json"
	"net/http"

	"dynamic-user-segmentation-service/internal/models"
	"github.com/gorilla/mux"
)

// GetSegmentsHandler godoc
//
//	@Summary		Получить семгенты
//	@Description	Получить все сегменты из БД
//	@ID				get-all-segments
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Success		200	    {array} 	models.Segment			"Все сегменты успешно получены"
//	@Failure		500	    {object}	models.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments [get]
func GetSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	segments, err := segmentRepository.GetAllSegments()

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех сегментов",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	// TODO: некрасиво, стоит переделать
	if segments == nil {
		segments = []*models.Segment{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(segments)
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
//	@Success		200	    {object} 	models.Segment			"Сегмент с данным названием успешно получен"
//	@Failure		400		{object}	models.ErrorDto			"Некорректные входные данные"
//	@Failure		404		{object}	models.ErrorDto			"Сегмент с данным названием не найден"
//	@Failure		500	    {object}	models.ErrorDto			"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments/{slug} [get]
func GetSegmentBySlugHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	segment, err := segmentRepository.GetSegmentBySlug(slug)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе сегмента по названию",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(segment)
}

// CreateSegmentHandler godoc
//
//		@Summary		Добавить сегмент
//		@Description	Добавить сегмент в БД
//		@ID				create-segment
//		@Tags			segments
//		@Accept			json
//		@Produce		json
//	 	@Param			Информация о сегменте 	body	models.CreateOrUpdateSegmentDto	    true	"Информация о добавляемом сегменте"
//		@Success		201		{object}	models.CreateSegmentResponseDto		"Сегмент успешно создан"
//		@Failure		400		{object}	models.ErrorDto						"Некорректные входные данные"
//		@Failure		500	    {object}	models.ErrorDto						"Возникла внутренняя ошибка сервера"
//		@Router			/api/v1/segments [post]
func CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	var segment models.CreateOrUpdateSegmentDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	slug, err := segmentRepository.CreateSegment(segment.Slug, segment.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании сегмента",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	createResponseDto := models.CreateSegmentResponseDto{
		Slug: slug,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createResponseDto)
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
//	 	@Param			Информация о сегменте	body	models.CreateOrUpdateSegmentDto	    true	"Информация о добавляемом сегменте"
//		@Success		200		{object}	models.UpdateSegmentResponseDto		"Сегмент с данным названием успешно обновлен"
//		@Failure		400		{object}	models.ErrorDto						"Некорректные входные данные"
//		@Failure		404		{object}	models.ErrorDto						"Сегмент с данным названием не найден"
//		@Failure		500	    {object}	models.ErrorDto						"Возникла внутренняя ошибка сервреа"
//		@Router			/api/v1/segments/{slug} [put]
func UpdateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	var segment models.CreateOrUpdateSegmentDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	updated, err := segmentRepository.UpdateSegment(slug, segment.Description)
	if err != nil {
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при обновлении сегмента",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	updateResponseDto := models.UpdateSegmentResponseDto{
		Id:          updated.Id,
		Slug:        updated.Slug,
		Description: updated.Description,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateResponseDto)
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
//	@Failure		400		{object}	models.ErrorDto						"Некорректные входные данные"
//	@Failure		404		{object}	models.ErrorDto						"Сегмент с данным названием не найден"
//	@Failure		500	    {object}	models.ErrorDto						"Возникла внутренняя ошибка сервера"
//	@Router			/api/v1/segments/{slug} [delete]
func DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	_, err := segmentRepository.DeleteSegment(slug)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при удалении сегмента",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
