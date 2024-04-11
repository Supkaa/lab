package entities

type Role struct {
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
}
