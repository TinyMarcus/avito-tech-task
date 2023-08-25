package models

import "database/sql"

type Segment struct {
	Id          int    `json:"id,omitempty"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
}

type UserSegment struct {
	UserId       int            `json:"user_id"`
	Slug         string         `json:"slug"`
	DeadlineDate sql.NullString `json:"deadline_date,omitempty"`
}

type CreateOrUpdateSegmentDto struct {
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
}

type SegmentWithDeadlineDate struct {
	Slug         string `json:"slug"`
	DeadlineDate string `json:"deadline_date,omitempty"`
}

type ChangeUserSegmentsDto struct {
	AddToUser    []SegmentWithDeadlineDate `json:"add_to_user"`
	TakeFromUser []string                  `json:"take_from_user"`
}

type UsersActiveSegments struct {
	UserId   int                       `json:"user_id"`
	Segments []SegmentWithDeadlineDate `json:"segments"`
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
