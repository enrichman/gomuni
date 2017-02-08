package gomuni

import (
	"errors"
	"math"
	"unicode"
)

const (
	k0 float64 = 0.9996
	e  float64 = 0.00669438
	r          = 6378137
)

var e2 = e * e
var e3 = e2 * e
var eP2 = e / (1.0 - e)

var sqrtE = math.Sqrt(1 - e)

var _e = (1 - sqrtE) / (1 + sqrtE)
var _e2 = _e * _e
var _e3 = _e2 * _e
var _e4 = _e3 * _e
var _e5 = _e4 * _e

var m1 = (1 - e/4 - 3*e2/64 - 5*e3/256)
var m2 = (3*e/8 + 3*e2/32 + 45*e3/1024)
var m3 = (15*e2/256 + 45*e3/1024)
var m4 = (35 * e3 / 3072)

var p2 = (3./2*_e - 27./32*_e3 + 269./512*_e5)
var p3 = (21./16*_e2 - 55./32*_e4)
var p4 = (151./96*_e3 - 417./128*_e5)
var p5 = (1097. / 512 * _e4)

const x = math.Pi / 180

func rad(d float64) float64 { return d * x }
func deg(r float64) float64 { return r / x }

type LatLng struct {
	lat float64
	lng float64
}

func toLatLon(easting, northing float64, zoneNumber int, zoneLetter string) (latlng *LatLng, err error) {

	x := easting - 500000
	y := northing

	var northernValue bool
	if zoneLetter != "" {
		zoneLetter := unicode.ToUpper(rune(zoneLetter[0]))
		if !('C' <= zoneLetter && zoneLetter <= 'X') || zoneLetter == 'I' || zoneLetter == 'O' {
			err := errors.New("zone letter out of range (must be between C and X)")
			return nil, err
		}
		northernValue = (zoneLetter >= 'N')
	}

	if !northernValue {
		y -= 10000000
	}

	m := y / k0
	mu := m / (r * m1)

	pRad := (mu +
		p2*math.Sin(2*mu) +
		p3*math.Sin(4*mu) +
		p4*math.Sin(6*mu) +
		p5*math.Sin(8*mu))

	pSin := math.Sin(pRad)
	pSin2 := pSin * pSin

	pCos := math.Cos(pRad)

	pTan := pSin / pCos
	pTan2 := pTan * pTan
	pTan4 := pTan2 * pTan2

	epSin := 1 - e*pSin2
	epSinSqrt := math.Sqrt(1 - e*pSin2)

	n := r / epSinSqrt
	rad := (1 - e) / epSin

	c := _e * pCos * pCos
	c2 := c * c

	d := x / (n * k0)
	d2 := d * d
	d3 := d2 * d
	d4 := d3 * d
	d5 := d4 * d
	d6 := d5 * d

	latitude := (pRad - (pTan/rad)*
		(d2/2-
			d4/24*(5+3*pTan2+10*c-4*c2-9*eP2)) +
		d6/720*(61+90*pTan2+298*c+45*pTan4-252*eP2-3*c2))

	longitude := (d -
		d3/6*(1+2*pTan2+c) +
		d5/120*(5-2*c+28*pTan2-3*c2+8*eP2+24*pTan4)) / pCos

	return &LatLng{deg(latitude), deg(longitude) + float64(zoneNumberToCentralLongitude(zoneNumber))}, nil
}

func zoneNumberToCentralLongitude(zoneNumber int) int {
	return (zoneNumber-1)*6 - 180 + 3
}
