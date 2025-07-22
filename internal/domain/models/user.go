package models

type User struct {
	Id           string
	Login        string
	PasswordHash string
	Email        string
	Name         string
	Surname      string
}
