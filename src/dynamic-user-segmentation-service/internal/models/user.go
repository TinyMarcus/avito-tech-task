package models

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateUserDto struct {
	Name string `json:"name"`
}
