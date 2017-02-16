package gomuni

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	shp "github.com/jonas-p/go-shp"
)

func initCountry() *Country {
	_ = godotenv.Load()
	rand.Seed(time.Now().Unix())

	regionFolder := os.Getenv("REGION_FOLDER")
	cityFolder := os.Getenv("CITY_FOLDER")
	townFolder := os.Getenv("TOWN_FOLDER")

	return Load(regionFolder, cityFolder, townFolder)
}

func getRandomPointsWithinRegions(regions []*Region, numOfPoints int) (points []Point) {

	points = make([]Point, numOfPoints)

	for i := 0; i < numOfPoints; i++ {
		reg := regions[rand.Int31n(int32(len(regions)))]
		lat := reg.BBox.MinX + (rand.Float64() * (reg.BBox.MaxX - reg.BBox.MinX))
		lng := reg.BBox.MinY + (rand.Float64() * (reg.BBox.MaxY - reg.BBox.MinY))

		points = append(points, Point{lat, lng})
	}

	return points
}

func getRandomPointsWithinBBoxes(boxes []*shp.Box, numOfPoints int) (points []Point) {

	points = make([]Point, numOfPoints)

	for i := 0; i < numOfPoints; i++ {
		box := boxes[rand.Int31n(int32(len(boxes)))]
		lat := box.MinX + (rand.Float64() * (box.MaxX - box.MinX))
		lng := box.MinY + (rand.Float64() * (box.MaxY - box.MinY))

		points = append(points, Point{lat, lng})
	}

	return points
}

var result interface{}

func Benchmark_FindTownByPoint(b *testing.B) {
	country := initCountry()

	boxes := make([]*shp.Box, 0)
	for _, r := range country.Regions {
		boxes = append(boxes, &r.BBox)
	}
	points := getRandomPointsWithinBBoxes(boxes, b.N)

	b.ResetTimer()

	var town *Town
	for i := 0; i < b.N; i++ {
		town = country.FindTownByPoint(points[i])
	}

	result = town
}

func Benchmark_GetRegionsByPoint(b *testing.B) {
	country := initCountry()

	boxes := make([]*shp.Box, 0)
	for _, r := range country.Regions {
		boxes = append(boxes, &r.BBox)
	}
	points := getRandomPointsWithinBBoxes(boxes, b.N)

	b.ResetTimer()

	var region []*Region
	for i := 0; i < b.N; i++ {
		region = country.GetRegionsByPoint(points[i])
	}

	result = region
}

func Benchmark_GetCitiesByPoint(b *testing.B) {
	country := initCountry()

	reg := country.Regions[rand.Int31n(int32(len(country.Regions)))]
	boxes := make([]*shp.Box, 0)
	for _, c := range reg.Cities {
		boxes = append(boxes, &c.BBox)
	}
	points := getRandomPointsWithinBBoxes(boxes, b.N)

	b.ResetTimer()

	var city []*City
	for i := 0; i < b.N; i++ {
		city = reg.GetCitiesByPoint(points[i])
	}

	result = city
}

func Benchmark_GetTownsByPoint(b *testing.B) {
	country := initCountry()

	reg := country.Regions[rand.Int31n(int32(len(country.Regions)))]
	city := reg.Cities[rand.Int31n(int32(len(reg.Cities)))]
	boxes := make([]*shp.Box, 0)
	for _, t := range city.Towns {
		boxes = append(boxes, &t.BBox)
	}
	points := getRandomPointsWithinBBoxes(boxes, b.N)

	b.ResetTimer()

	var town []*Town
	for i := 0; i < b.N; i++ {
		town = city.GetTownsByPoint(points[i])
	}

	result = town
}
