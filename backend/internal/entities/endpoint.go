package entities

import (
	"github.com/google/uuid"
)

type Endpoint struct {
	ID     uuid.UUID `json:"id"`
	URL    string    `json:"url"`
	Method string    `json:"method"`
}
