package models

import "errors"

type ExoplanetType string

const (
	GasGiant   ExoplanetType = "GasGiant"
	Terretrial ExoplanetType = "Terretrial"
)

type Exoplanet struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Distance    float64       `json:"distance"`
	Radius      float64       `json:"radius"`
	Type        ExoplanetType `json:"type"`
	Mass        *float64      `json:"mass,omitempty"`
}

type ExoplanetStore struct {
	Exoplanets map[string]Exoplanet
}

func NewExoPlanetStore() *ExoplanetStore {
	return &ExoplanetStore{
		Exoplanets: make(map[string]Exoplanet),
	}
}

var (
	ErrNotFound = errors.New("Exoplanet not found")
	ErrInvalid  = errors.New("invalid exoplanet data")
)
