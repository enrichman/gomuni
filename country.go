package gomuni

import (
	"io/ioutil"
	"log"
	"strings"

	"errors"

	"github.com/dhconnelly/rtreego"
	shp "github.com/jonas-p/go-shp"
	geo "github.com/kellydunn/golang-geo"
)

//Country represent the Italy with its regions
type Country struct {
	Regions []*Region `json:"regions,omitempty"`

	regionsTree *rtreego.Rtree
	regionsMap  map[string]*Region
}

//RegionsGetter can be used to retrive a region from its ID or from a geolocation point
type RegionsGetter interface {
	GetRegionById(ID string) *Region
	GetRegionsByPoint(lat, lng float32) []*Region
}

//Load all the country with the Regions, Cities and Towns
func Load(regionFolder, cityFolder, townFolder string) *Country {
	country := loadCountryWithRegions(regionFolder)
	country.loadRegionsWithCities(cityFolder)
	country.loadCitiesWithTowns(townFolder)
	return country
}

//GetRegionById return a Region with the provided ID
func (c *Country) GetRegionById(ID string) *Region {
	return c.regionsMap[ID]
}

//GetRegionsByPoint return the Regions with the bounding box over the passed geolocation
func (c *Country) GetRegionsByPoint(lat, lng float64) []*Region {
	location := rtreego.Point{lat, lng}
	results := c.regionsTree.SearchIntersect(location.ToRect(0.01))

	regions := make([]*Region, 0)
	for _, s := range results {
		r := s.(*Region)
		regions = append(regions, r)
	}

	return regions
}

func loadCountryWithRegions(folder string) *Country {
	regions := make([]*Region, 0)
	regionsTree := rtreego.NewTree(2, 25, 50)
	regionsMap := make(map[string]*Region)

	loaded := false
	files, _ := ioutil.ReadDir(folder)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".shp") {
			reader, err := shp.Open(folder + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			for reader.Next() {
				n, s := reader.Shape()
				switch s.(type) {
				case *shp.Polygon:
					codReg := reader.ReadAttribute(n, 0)
					nameReg := reader.ReadAttribute(n, 1)

					reg := &Region{
						ID:         codReg,
						Name:       nameReg,
						Cities:     make([]*City, 0),
						citiesTree: rtreego.NewTree(2, 25, 50),
						citiesMap:  make(map[string]*City),
					}

					p := s.(*shp.Polygon)

					// load bounding box
					minLatLng, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxLatLng, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					reg.BBox = shp.Box{
						MinX: minLatLng.lat,
						MinY: minLatLng.lng,
						MaxX: maxLatLng.lat,
						MaxY: maxLatLng.lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						latlng, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(latlng.lat, latlng.lng))
					}
					reg.polygon = geo.NewPolygon(points)

					regions = append(regions, reg)
					regionsMap[reg.ID] = reg
					regionsTree.Insert(reg)
				}
			}

			loaded = true
		}
	}

	if !loaded {
		err := errors.New("Regions not loaded!")
		panic(err)
	}

	return &Country{regions, regionsTree, regionsMap}
}

func (c *Country) loadRegionsWithCities(folder string) {
	files, _ := ioutil.ReadDir(folder)

	loaded := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".shp") {
			reader, err := shp.Open(folder + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			for reader.Next() {
				n, s := reader.Shape()
				switch s.(type) {
				case *shp.Polygon:
					regID := reader.ReadAttribute(n, 0)
					cityID := reader.ReadAttribute(n, 2)
					cityName := reader.ReadAttribute(n, 3)
					cityShortname := reader.ReadAttribute(n, 5)
					maincityFlag := reader.ReadAttribute(n, 6)

					city := &City{
						RegionID:  regID,
						ID:        cityID,
						Name:      cityName,
						Shortname: cityShortname,
						Maincity:  (maincityFlag == "1"),
						Towns:     make([]*Town, 0),
						townsTree: rtreego.NewTree(2, 25, 350),
						townsMap:  make(map[string]*Town),
					}

					p := s.(*shp.Polygon)

					// load bounding box
					minLatLng, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxLatLng, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					city.BBox = shp.Box{
						MinX: minLatLng.lat,
						MinY: minLatLng.lng,
						MaxX: maxLatLng.lat,
						MaxY: maxLatLng.lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						latlng, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(latlng.lat, latlng.lng))
					}
					city.polygon = geo.NewPolygon(points)

					region := c.GetRegionById(regID)
					region.addCity(city)
				}
			}

			loaded = true
		}
	}

	if !loaded {
		err := errors.New("Cities not loaded!")
		panic(err)
	}
}

func (c *Country) loadCitiesWithTowns(folder string) {
	files, _ := ioutil.ReadDir(folder)

	loaded := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".shp") {
			reader, err := shp.Open(folder + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			for reader.Next() {
				n, s := reader.Shape()
				switch s.(type) {
				case *shp.Polygon:
					regID := reader.ReadAttribute(n, 0)
					cityID := reader.ReadAttribute(n, 2)
					townID := reader.ReadAttribute(n, 3)
					townName := reader.ReadAttribute(n, 4)

					town := &Town{
						ID:       townID,
						RegionID: regID,
						CityID:   cityID,
						Name:     townName,
					}

					p := s.(*shp.Polygon)

					// load bounding box
					minLatLng, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxLatLng, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					town.BBox = shp.Box{
						MinX: minLatLng.lat,
						MinY: minLatLng.lng,
						MaxX: maxLatLng.lat,
						MaxY: maxLatLng.lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						latlng, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(latlng.lat, latlng.lng))
					}
					town.polygon = geo.NewPolygon(points)

					region := c.GetRegionById(regID)
					city := region.GetCityById(cityID)
					city.addTown(town)
				}
			}

			loaded = true
		}
	}

	if !loaded {
		err := errors.New("Towns not loaded!")
		panic(err)
	}
}
