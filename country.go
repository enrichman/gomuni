package gomuni

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/dhconnelly/rtreego"
	shp "github.com/jonas-p/go-shp"
)

type Country struct {
	Regions []*Region

	regionsTree *rtreego.Rtree
	regionsMap  map[string]*Region
}

type RegionsGetter interface {
	GetRegionById(ID string) *Region
	GetRegionsByPoint(lat, lng float32) []*Region
}

func Load(regionFolder, cityFolder, townFolder string) *Country {
	country := loadCountryWithRegions(regionFolder)
	country.loadRegionsWithCities(cityFolder)
	country.loadCitiesWithTowns(townFolder)
	return country
}

func (c *Country) GetRegionById(ID string) *Region {
	return c.regionsMap[ID]
}

func (c *Country) GetRegionsByPoint(lat, lng float32) []*Region {
	return nil //c.GetRegionsByPoint(lat, lng)
}

func loadCountryWithRegions(folder string) *Country {
	regions := make([]*Region, 0)
	regionsTree := rtreego.NewTree(2, 25, 50)
	regionsMap := make(map[string]*Region)

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
						Code:       codReg,
						Name:       nameReg,
						Cities:     make([]*City, 0),
						citiesTree: rtreego.NewTree(2, 25, 50),
						citiesMap:  make(map[string]*City),
					}
					regions = append(regions, reg)
					regionsMap[reg.Code] = reg
					regionsTree.Insert(reg)

					//p := s.(*shp.Polygon)
				}
			}
		}
	}

	return &Country{regions, regionsTree, regionsMap}
}

func (c *Country) loadRegionsWithCities(folder string) {
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

					region := c.GetRegionById(regID)
					region.addCity(city)

					//p := s.(*shp.Polygon)
				}
			}
		}
	}
}

func (c *Country) loadCitiesWithTowns(folder string) {
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

					region := c.GetRegionById(regID)
					city := region.GetCityById(cityID)
					city.addTown(town)

					//p := s.(*shp.Polygon)
				}
			}
		}
	}
}
