package gomuni

import "github.com/dhconnelly/rtreego"

type Town struct {
	ID       string
	RegionID string
	CityID   string
	Name     string
}

func (t *Town) Bounds() *rtreego.Rect {
	return nil
}
