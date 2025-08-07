package dto

import "github.com/google/uuid"

type AddBooksToCustomerRequest struct {
	BookIDs []uuid.UUID `json:"bookIDs"`
}
