package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/vaishnavitnaik/models"
	"github.com/vaishnavitnaik/utils"
)

var store = models.NewExoPlanetStore()

func respondwithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"Error": message})
}

func validateExoPlanet(exoplanet *models.Exoplanet) error {
	if exoplanet.Name == "" || exoplanet.Description == "" {
		return models.ErrInvalid
	}
	if exoplanet.Distance <= 10 || exoplanet.Distance >= 1000 {
		return models.ErrInvalid
	}
	if exoplanet.Radius <= -0.1 || exoplanet.Radius >= 10 {
		return models.ErrInvalid
	}
	if exoplanet.Type != models.GasGiant && exoplanet.Type != models.Terretrial {
		return models.ErrInvalid
	}
	return nil
}

func AddExoPlanet(w http.ResponseWriter, r *http.Request) {
	var exoplanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		respondwithError(w, http.StatusBadRequest, "Invalid data")
		return
	}
	if err := validateExoPlanet(&exoplanet); err != nil {
		respondwithError(w, http.StatusBadRequest, err.Error())
		return
	}
	exoplanet.ID = uuid.NewString()
	store.Exoplanets[exoplanet.ID] = exoplanet
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exoplanet)
}

func ListExoPlanet(w http.ResponseWriter, r *http.Request) {
	exoplanets := make([]models.Exoplanet, 0, len(store.Exoplanets))

	for _, planet := range store.Exoplanets {
		exoplanets = append(exoplanets, planet)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exoplanets)
}

func GetExoPlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exoplanet, ok := store.Exoplanets[id]
	if !ok {
		respondwithError(w, http.StatusNotFound, "invalid id or exoplanet")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exoplanet)
}

func UpdateExoPlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var reqPlanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&reqPlanet); err != nil {
		respondwithError(w, http.StatusBadRequest, "invalid request payload")
	}
	if err := validateExoPlanet(&reqPlanet); err != nil {
		respondwithError(w, http.StatusBadRequest, "Invalid exoplanet data")
	}
	_, ok := store.Exoplanets[id]
	if !ok {
		respondwithError(w, http.StatusNotFound, "invalid id or exoplanet")
	}
	reqPlanet.ID = id
	store.Exoplanets[id] = reqPlanet
	w.Header().Set("Content-Type", "applicatin/json")
	json.NewEncoder(w).Encode(reqPlanet)
}

func DeleteExoPlanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var planet models.Exoplanet
	if _, ok := store.Exoplanets[id]; !ok {
		respondwithError(w, http.StatusNotFound, models.ErrNotFound.Error())
		return
	}
	delete(store.Exoplanets, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Deleted successfully")
	json.NewEncoder(w).Encode(planet)

}
func FuelEstimation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	planet, ok := store.Exoplanets[id]
	if !ok {
		respondwithError(w, http.StatusNotFound, models.ErrNotFound.Error())
		return
	}
	crewCapacityStr := r.URL.Query().Get("crew")
	crewCapacity, err := strconv.Atoi(crewCapacityStr)
	if err != nil || crewCapacity <= 0 {
		respondwithError(w, http.StatusBadRequest, "Invalid crew capacity")
		return
	}
	fuelcost, err := utils.CalculateFuel(planet, crewCapacity)
	if err != nil {
		respondwithError(w, http.StatusBadRequest, err.Error())
		return
	}
	response := map[string]float32{"FuelCost": float32(fuelcost)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
