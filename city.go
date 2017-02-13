package gomuni

import (
	"github.com/dhconnelly/rtreego"
	shp "github.com/jonas-p/go-shp"
	geo "github.com/kellydunn/golang-geo"
)

type City struct {
	ID        string  `json:"id,omitempty"`
	RegionID  string  `json:"region_id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Shortname string  `json:"shortname,omitempty"`
	Maincity  bool    `json:"maincity,omitempty"`
	Towns     []*Town `json:"towns,omitempty"`

	BBox      shp.Box `json:"bbox,omitempty"`
	polygon   *geo.Polygon
	townsTree *rtreego.Rtree
	townsMap  map[string]*Town
}

type TownGetter interface {
	GetTownById(ID string) *Town
	GetTownByPoint(lat, lng float32) []*Town
}

func (c *City) Bounds() *rtreego.Rect {
	p1 := rtreego.Point{c.BBox.MinX, c.BBox.MinY}
	r1, _ := rtreego.NewRect(p1, []float64{c.BBox.MaxX - c.BBox.MinX, c.BBox.MaxY - c.BBox.MinY})
	return r1
}

func (c *City) GetTownById(ID string) *Town {
	return c.townsMap[ID]
}

func (c *City) GetTownByPoint(lat, lng float64) []*Town {
	location := rtreego.Point{lat, lng}
	results := c.townsTree.SearchIntersect(location.ToRect(0.01))

	cities := make([]*Town, 0)
	for _, s := range results {
		r := s.(*Town)
		cities = append(cities, r)
	}

	return cities
}

func (c *City) addTown(town *Town) {
	c.Towns = append(c.Towns, town)
	c.townsMap[town.ID] = town
	c.townsTree.Insert(town)
}
