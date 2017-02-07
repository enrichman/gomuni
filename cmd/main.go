package main

import (
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/enrichman/gomuni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	regionFolder := os.Getenv("REGION_FOLDER")
	cityFolder := os.Getenv("CITY_FOLDER")
	townFolder := os.Getenv("TOWN_FOLDER")

	country := gomuni.Load(regionFolder, cityFolder, townFolder)
	service := service{country}

	router := mux.NewRouter()
	router.HandleFunc("/country", service.countryHandler).Methods("GET")
	router.HandleFunc("/country/regions", service.regionsHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}", service.regionIDHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities", service.regionCitiesHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}", service.regionCityIDHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}/towns", service.townsHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}/towns/{town_id}", service.regionCityTownIDHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type service struct {
	country *gomuni.Country
}

func (s *service) countryHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(s.country)
	w.Write(b)
}

func (s *service) regionsHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(s.country.Regions)
	w.Write(b)
}

func (s *service) regionIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	b, _ := json.Marshal(s.country.GetRegionById(vars["region_id"]))
	w.Write(b)
}

func (s *service) regionCitiesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionById(vars["region_id"])
	b, _ := json.Marshal(region.Cities)
	w.Write(b)
}

func (s *service) regionCityIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionById(vars["region_id"])
	city := region.GetCityById(vars["city_id"])
	b, _ := json.Marshal(city)
	w.Write(b)
}

func (s *service) townsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionById(vars["region_id"])
	city := region.GetCityById(vars["city_id"])
	b, _ := json.Marshal(city.Towns)
	w.Write(b)
}

func (s *service) regionCityTownIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionById(vars["region_id"])
	city := region.GetCityById(vars["city_id"])
	town := city.GetTownById(vars["town_id"])
	b, _ := json.Marshal(town)
	w.Write(b)
}
