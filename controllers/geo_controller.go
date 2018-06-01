package controllers

import (
	"net/http"

	"github.com/b0ralgin/geo-ip-finder/services"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
)

type GeoIpController struct {
	db           *bolt.DB
	geoIpService *services.GeoIpService
}

func NewGeoController(db *bolt.DB, service *services.GeoIpService) *GeoIpController {
	return &GeoIpController{
		db:           db,
		geoIpService: service,
	}
}

func (g *GeoIpController) GetCountryByIp(c echo.Context) error {
	ip := c.RealIP()
	res, err := g.geoIpService.MakeRequest(ip)
	if err != nil {
		return c.String(http.StatusBadGateway, err.Error())
	}
	return c.String(http.StatusOK, res)
}
