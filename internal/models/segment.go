package models

import "database/sql"

type Segment struct {
	Id          int
	Slug        string
	Description string
}

type UserSegment struct {
	UserId       int
	Slug         string
	DeadlineDate sql.NullString
}
