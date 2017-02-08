package main

import (
	"log"
	"net/http"
	"os"

	"encoding/json"

	"strconv"

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
	Region *gomuni.Region
	City   *gomuni.City
	Town   *gomuni.Town
}

type service struct {
	country *gomuni.Country
}

func (s *service) searchHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	lat, okLat := vals["lat"]
	lng, okLng := vals["lng"]

	containedTowns := make([]*gomuni.Town, 0)

	if okLat && okLng {
		latFloat, _ := strconv.ParseFloat(lat[0], 64)
		lngFloat, _ := strconv.ParseFloat(lng[0], 64)
		regions := s.country.GetRegionsByPoint(latFloat, lngFloat)

		allTowns := make([]*gomuni.Town, 0)
		for _, r := range regions {
			cities := r.GetCityByPoint(latFloat, lngFloat)
			for _, c := range cities {
				towns := c.GetTownByPoint(latFloat, lngFloat)
				allTowns = append(allTowns, towns...)
			}
		}

		for _, t := range allTowns {
			if t.Contains(latFloat, lngFloat) {
				containedTowns = append(containedTowns, t)
			}
		}
	}

	b, _ := json.Marshal(containedTowns)
	w.Write(b)
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