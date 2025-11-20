package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Surname     string    `json:"surname" db:"surname"`
	Contacts    []string  `json:"contacts" db:"contacts"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Description string    `json:"description" db:"description"`
}
