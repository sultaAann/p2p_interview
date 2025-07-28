package models

import "p2p_interview/internal/domain/errors"

type User struct {
	Id           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"password,omitempty"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Surname      string `json:"surname,omitempty"`
}

func (user *User) Validate() error {
	if user.Login == "" {
		return &errors.ValidationError{Field: "login", Message: "cannot be empty"}
	}
	if user.PasswordHash == "" {
		return &errors.ValidationError{Field: "password", Message: "cannot be empty"}
	}
	if user.Email == "" {
		return &errors.ValidationError{Field: "email", Message: "cannot be empty"}
	}
	if user.Name == "" {
		return &errors.ValidationError{Field: "name", Message: "cannot be empty"}
	}
	// if user.Surname == "" {
	// 	return &errors.ValidationError{Field: "surname", Message: "cannot be empty"}
	// }
	return nil
}
