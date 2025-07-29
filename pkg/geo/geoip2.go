package geo

import (
	_ "embed"

	"github.com/oschwald/geoip2-golang"
)

//go:embed data/GeoLite2-City.mmdb
var geoDB []byte

var G *geoip2.Reader

func init() {
	var err error
	G, err = geoip2.FromBytes(geoDB)
	if err != nil {
		panic(err)
	}
}
