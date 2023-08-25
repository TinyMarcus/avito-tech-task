package errors

import "errors"

var (
	RecordNotFound       = errors.New("Record was not found")
	InvalidRequest       = errors.New("Invalid rerquest")
	DatabaseWritingError = errors.New("Error while writing to DB")
	DatabaseReadingError = errors.New("Error while reading from DB")
	RecordAlreadyExists  = errors.New("Record with this data already exists")
	UnknownError         = errors.New("Unknown error was happened")
)
