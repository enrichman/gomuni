package gomuni

import "github.com/dhconnelly/rtreego"

type City struct {
	ID        string
	RegionID  string
	Name      string
	Shortname string
	Maincity  bool
	Towns     []*Town

	townsTree *rtreego.Rtree
	townsMap  map[string]*Town
}

type TownGetter interface {
	GetTownById(ID string) *Town
	GetTownByPoint(lat, lng float32) []*Town
}

func (c *City) Bounds() *rtreego.Rect {
	return nil
}

func (c *City) GetTownById(ID string) *Town {
	return c.townsMap[ID]
}

func (c *City) GetTownByPoint(lat, lng float32) []*Town {
	return nil //c.GetRegionsByPoint(lat, lng)
}

func (c *City) addTown(town *Town) {
	c.Towns = append(c.Towns, town)
	c.townsMap[town.ID] = town
	c.townsTree.Insert(town)
}
