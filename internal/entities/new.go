package entities

import (
	"github.com/google/uuid"
	"time"
)

type New struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Views     int       `json:"views"`
}
