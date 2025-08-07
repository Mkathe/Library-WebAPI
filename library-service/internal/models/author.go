package models

import "github.com/google/uuid"

type AuthorModel struct {
	AuthorID       uuid.UUID `json:"author_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ThirdName      string    `json:"third_name"`
	Nickname       string    `json:"nickname"`
	Specialization string    `json:"specialization"`
}
