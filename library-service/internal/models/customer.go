package models

import "github.com/google/uuid"

type CustomerModel struct {
	CustomerID uuid.UUID `json:"customer_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ThirdName  string    `json:"third_name"`
}
