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

//Point represent a geolocation point with latitude and longitude
type Point struct {
	Lat float64
	Lng float64
}

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

//GetRegionByID returns the Region with the provided ID
func (c *Country) GetRegionByID(ID string) *Region {
	return c.regionsMap[ID]
}

//GetRegionsByPoint returns the Regions having their bounding box over the provided geolocation point
func (c *Country) GetRegionsByPoint(point Point) []*Region {
	location := rtreego.Point{point.Lat, point.Lng}
	results := c.regionsTree.SearchIntersect(location.ToRect(0.01))

	regions := make([]*Region, 0)
	for _, s := range results {
		r := s.(*Region)
		regions = append(regions, r)
	}

	return regions
}

//FindTownByPoint return the closest Town from the  Point
func (c *Country) FindTownByPoint(point Point) *Town {
	containedTowns := make([]*Town, 0)

	regions := c.GetRegionsByPoint(point)

	allTowns := make([]*Town, 0)
	for _, r := range regions {
		cities := r.GetCitiesByPoint(point)
		for _, c := range cities {
			towns := c.GetTownsByPoint(point)
			allTowns = append(allTowns, towns...)
		}
	}

	for _, t := range allTowns {
		if t.Contains(point) {
			containedTowns = append(containedTowns, t)
		}
	}

	if len(containedTowns) > 0 {
		return containedTowns[0]
	}

	return nil
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
					minPoint, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxPoint, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					reg.BBox = shp.Box{
						MinX: minPoint.Lat,
						MinY: minPoint.Lng,
						MaxX: maxPoint.Lat,
						MaxY: maxPoint.Lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						point, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(point.Lat, point.Lng))
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
					minPoint, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxPoint, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					city.BBox = shp.Box{
						MinX: minPoint.Lat,
						MinY: minPoint.Lng,
						MaxX: maxPoint.Lat,
						MaxY: maxPoint.Lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						point, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(point.Lat, point.Lng))
					}
					city.polygon = geo.NewPolygon(points)

					region := c.GetRegionByID(regID)
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
						ID:       buildIstatID(townID),
						RegionID: regID,
						CityID:   cityID,
						Name:     townName,
					}

					p := s.(*shp.Polygon)

					// load bounding box
					minPoint, _ := toLatLon(p.Box.MinX, p.Box.MinY, 32, "N")
					maxPoint, _ := toLatLon(p.Box.MaxX, p.Box.MaxY, 32, "N")
					town.BBox = shp.Box{
						MinX: minPoint.Lat,
						MinY: minPoint.Lng,
						MaxX: maxPoint.Lat,
						MaxY: maxPoint.Lng,
					}

					// load polygon
					points := make([]*geo.Point, 0)
					for _, point := range p.Points {
						point, _ := toLatLon(point.X, point.Y, 32, "N")
						points = append(points, geo.NewPoint(point.Lat, point.Lng))
					}
					town.polygon = geo.NewPolygon(points)

					region := c.GetRegionByID(regID)
					city := region.GetCityByID(cityID)
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

func buildIstatID(id string) string {
	for len(id) < 6 {
		id = "0" + id
	}
	return id
}
