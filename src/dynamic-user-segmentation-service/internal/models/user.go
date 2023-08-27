package models

// User model info
// @Description Информация о пользователе
type User struct {
	Id   int    `json:"id"`   // Идентификатор пользователя
	Name string `json:"name"` // Имя пользователя
}

// CreateUserDto model info
// @Description Информация о пользователе при создании
type CreateUserDto struct {
	Name string `json:"name"` // Имя пользователя
}

// CreateUserResponseDto model info
// @Description Информация о пользователе при создании
type CreateUserResponseDto struct {
	Id int `json:"id"` // Идентификатор пользователя
}
