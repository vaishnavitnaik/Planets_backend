package utils

import (
	"errors"
	"math"

	"github.com/vaishnavitnaik/models"
)

func CalculateFuel(exoplanet models.Exoplanet, crewCapacity int) (float64, error) {
	var g float64
	switch exoplanet.Type {
	case models.GasGiant:
		g = 0.5 / math.Pow(exoplanet.Radius, 2)
	case models.Terretrial:
		if exoplanet.Mass == nil {
			return 0, errors.New("Exoplanet cannot have mass 0")
		}
		g *= *exoplanet.Mass / math.Pow(exoplanet.Radius, 2)
	default:
		return 0, errors.New("unknown type")
	}
	fuelCost := exoplanet.Distance / math.Pow(g, 2) * float64(crewCapacity)
	return fuelCost, nil
}
