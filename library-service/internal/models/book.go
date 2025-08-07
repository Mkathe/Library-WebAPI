package models

import "github.com/google/uuid"

type BookModel struct {
	BookID uuid.UUID `json:"book_id"`
	Title  string    `json:"title"`
	Genre  string    `json:"genre"`
	ISBN   string    `json:"isbn"`
}
