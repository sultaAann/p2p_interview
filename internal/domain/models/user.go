package models

type User struct {
	Id           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"passwordhash"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
}
