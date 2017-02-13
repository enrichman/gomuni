package gomuni

import (
	"github.com/dhconnelly/rtreego"
	shp "github.com/jonas-p/go-shp"
	geo "github.com/kellydunn/golang-geo"
)

type Town struct {
	ID       string `json:"id,omitempty"`
	RegionID string `json:"region_id,omitempty"`
	CityID   string `json:"city_id,omitempty"`
	Name     string `json:"name,omitempty"`

	BBox    shp.Box `json:"bbox,omitempty"`
	polygon *geo.Polygon
}

func (t *Town) Bounds() *rtreego.Rect {
	p1 := rtreego.Point{t.BBox.MinX, t.BBox.MinY}
	r1, _ := rtreego.NewRect(p1, []float64{t.BBox.MaxX - t.BBox.MinX, t.BBox.MaxY - t.BBox.MinY})
	return r1
}

func (t *Town) Contains(lat, lng float64) bool {
	return t.polygon.Contains(geo.NewPoint(lat, lng))
}
