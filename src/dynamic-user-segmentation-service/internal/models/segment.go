package models

import "database/sql"

// Segment model info
// @Description Информация о сегменте
type Segment struct {
	Id          int    `json:"id,omitempty"`          // Идентификатор сегмента
	Slug        string `json:"slug"`                  // Название сегмента
	Description string `json:"description,omitempty"` // Описание сегмента
}

type UserSegment struct {
	UserId       int            `json:"user_id"`
	Slug         string         `json:"slug"`
	DeadlineDate sql.NullString `json:"deadline_date,omitempty"`
}

// CreateOrUpdateSegmentDto model info
// @Description Информация о сегменте при создании
type CreateOrUpdateSegmentDto struct {
	Slug        string `json:"slug"`                  // Название сегмента
	Description string `json:"description,omitempty"` // Описание сегмента
}

// SegmentWithDeadlineDate model info
// @Description Информация о сегментах с датой отключения пользователя от сегмента
type SegmentWithDeadlineDate struct {
	Slug         string `json:"slug"`                    // Название сегмента
	DeadlineDate string `json:"deadline_date,omitempty"` // Дата отключения пользователя от сегмента
}

// ChangeUserSegmentsDto model info
// @Description Информация о добавляемых и удаляемых сегментах пользователя
type ChangeUserSegmentsDto struct {
	AddToUser    []SegmentWithDeadlineDate `json:"add_to_user"`    // Сегменты, которые будут добавляться пользователю (с датами отключения)
	TakeFromUser []string                  `json:"take_from_user"` // Сегменты, которые будут удаляться у пользователя
}

// UsersActiveSegments model info
// @Description Информация об активных сегментах пользователя
type UsersActiveSegments struct {
	UserId   int                       `json:"user_id"`  // Идентификатор пользователя
	Segments []SegmentWithDeadlineDate `json:"segments"` // Список активных сегментов
}

func ConvertUserSegmentToUsersActiveSegments(userId int, userSegments []*UserSegment) *UsersActiveSegments {
	var segments []SegmentWithDeadlineDate

	for _, val := range userSegments {
		segmentWithDeadlineDate := &SegmentWithDeadlineDate{
			Slug:         val.Slug,
			DeadlineDate: val.DeadlineDate.String,
		}
		segments = append(segments, *segmentWithDeadlineDate)
	}

	return &UsersActiveSegments{
		UserId:   userId,
		Segments: segments,
	}
}
