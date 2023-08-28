package dto

import "github.com/TinyMarcus/avito-tech-task/internal/models"

// UserDto model info
// @Description Информация о пользователе
type UserDto struct {
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

func ConvertUserToUserDto(user *models.User) *UserDto {
	return &UserDto{
		Id:   user.Id,
		Name: user.Name,
	}
}
