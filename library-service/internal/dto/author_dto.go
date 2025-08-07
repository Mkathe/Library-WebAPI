package dto

import "github.com/google/uuid"

type AddBooksToAuthorRequest struct {
	BookIDs []uuid.UUID `json:"bookIDs"`
}
