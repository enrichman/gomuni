package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"encoding/json"

	"strconv"

	"github.com/enrichman/gofield"
	"github.com/enrichman/gomuni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println(".env file loaded")

	regionFolder := os.Getenv("REGION_FOLDER")
	cityFolder := os.Getenv("CITY_FOLDER")
	townFolder := os.Getenv("TOWN_FOLDER")

	log.Println("Loading folders:", regionFolder, cityFolder, townFolder)
	country := gomuni.Load(regionFolder, cityFolder, townFolder)
	log.Println("Country loaded")

	log.Println("Loading handlers")
	service := service{country}

	router := mux.NewRouter()
	router.HandleFunc("/search", service.searchHandler).Methods("GET")
	router.HandleFunc("/country", service.countryHandler).Methods("GET")
	router.HandleFunc("/country/regions", service.regionsHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}", service.regionIDHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities", service.regionCitiesHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}", service.regionCityIDHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}/towns", service.townsHandler).Methods("GET")
	router.HandleFunc("/country/regions/{region_id}/cities/{city_id}/towns/{town_id}", service.regionCityTownIDHandler).Methods("GET")

	log.Println("Ready")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type response struct {
	Region *gomuni.Region `json:"region,omitempty"`
	City   *gomuni.City   `json:"city,omitempty"`
	Town   *gomuni.Town   `json:"town,omitempty"`
}

type service struct {
	country *gomuni.Country
}

func (s *service) searchHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()

	var lat string
	var lng string

	latlng, okLatLng := vals["latlng"]
	if okLatLng {
		latlng = strings.Split(latlng[0], ",")
	}

	latArr, okLat := vals["lat"]
	lngArr, okLng := vals["lng"]

	if okLatLng && len(latlng) > 1 {
		lat = latlng[0]
		lng = latlng[1]
	} else if okLat && okLng {
		lat = latArr[0]
		lng = lngArr[0]
	}

	var town *gomuni.Town
	if lat != "" && lng != "" {
		latFloat, _ := strconv.ParseFloat(lat, 64)
		lngFloat, _ := strconv.ParseFloat(lng, 64)
		point := gomuni.Point{Lat: latFloat, Lng: lngFloat}
		town = s.country.FindTownByPoint(point)
	}

	b, _ := json.Marshal(town)
	w.Write(b)
}

func (s *service) countryHandler(w http.ResponseWriter, r *http.Request) {
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(s.country, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) regionsHandler(w http.ResponseWriter, r *http.Request) {
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(s.country.Regions, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) regionIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(s.country.GetRegionByID(vars["region_id"]), fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) regionCitiesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionByID(vars["region_id"])
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(region.Cities, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) regionCityIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionByID(vars["region_id"])
	city := region.GetCityByID(vars["city_id"])
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(city, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) townsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionByID(vars["region_id"])
	city := region.GetCityByID(vars["city_id"])
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(city.Towns, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}

func (s *service) regionCityTownIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := s.country.GetRegionByID(vars["region_id"])
	city := region.GetCityByID(vars["city_id"])
	town := city.GetTownByID(vars["town_id"])
	fields := r.URL.Query().Get("fields")
	lightObj := gofield.Reduce(town, fields)
	b, _ := json.Marshal(lightObj)
	w.Write(b)
}
