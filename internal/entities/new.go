package entities

import "time"

type New struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Views     int       `json:"views"`
}
