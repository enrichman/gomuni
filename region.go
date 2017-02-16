package gomuni

import (
	"github.com/dhconnelly/rtreego"
	shp "github.com/jonas-p/go-shp"
	geo "github.com/kellydunn/golang-geo"
)

//Region represent an italian Region with its cities
type Region struct {
	ID     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Cities []*City `json:"cities,omitempty"`

	BBox       shp.Box `json:"bbox,omitempty"`
	polygon    *geo.Polygon
	citiesTree *rtreego.Rtree
	citiesMap  map[string]*City
}

//CityGetter can be used to retrive a city from its ID or from a geolocation point
type CityGetter interface {
	GetCityById(ID string) *City
	GetCityByPoint(lat, lng float32) []*City
}

//Bounds is used to implement the rtreego Spatial interface
func (r *Region) Bounds() *rtreego.Rect {
	p1 := rtreego.Point{r.BBox.MinX, r.BBox.MinY}
	r1, _ := rtreego.NewRect(p1, []float64{r.BBox.MaxX - r.BBox.MinX, r.BBox.MaxY - r.BBox.MinY})
	return r1
}

//GetCityByID returns the City with the provided ID
func (r *Region) GetCityByID(ID string) *City {
	return r.citiesMap[ID]
}

//GetCitiesByPoint returns the Cities having their bounding box over the provided geolocation point
func (r *Region) GetCitiesByPoint(point Point) []*City {
	location := rtreego.Point{point.Lat, point.Lng}
	results := r.citiesTree.SearchIntersect(location.ToRect(0.01))

	cities := make([]*City, 0)
	for _, s := range results {
		r := s.(*City)
		cities = append(cities, r)
	}

	return cities
}

func (r *Region) addCity(city *City) {
	r.Cities = append(r.Cities, city)
	r.citiesMap[city.ID] = city
	r.citiesTree.Insert(city)
}
