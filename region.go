package gomuni

import "github.com/dhconnelly/rtreego"

type Region struct {
	Code   string
	Name   string
	Cities []*City

	citiesTree *rtreego.Rtree
	citiesMap  map[string]*City
}

type CityGetter interface {
	GetCityById(ID string) *City
	GetCityByPoint(lat, lng float32) []*City
}

func (r *Region) Bounds() *rtreego.Rect {
	return nil
}

func (r *Region) GetCityById(ID string) *City {
	return r.citiesMap[ID]
}

func (r *Region) GetCityByPoint(lat, lng float32) []*City {
	return nil //c.GetRegionsByPoint(lat, lng)
}

func (r *Region) addCity(city *City) {
	r.Cities = append(r.Cities, city)
	r.citiesMap[city.ID] = city
	r.citiesTree.Insert(city)
}
