package models

import "p2p_interview/internal/domain/models"

type UserResponse struct {
	Id      string `json:"id"`
	Login   string `json:"login"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname,omitempty"`
}

func UserToUserResponse(user models.User) UserResponse {
	return UserResponse{Id: user.Id, Login: user.Login, Email: user.Email, Name: user.Name, Surname: user.Surname}
}
