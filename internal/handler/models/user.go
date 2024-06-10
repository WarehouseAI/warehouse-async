package models

type (
	CreateUserRequest struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Username  string `json:"username"`
		Hash      string `json:"hash"`
		Email     string `json:"email"`
	}
)
