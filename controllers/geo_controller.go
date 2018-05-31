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

func (g *GeoIpController) GetIp(c echo.Context) error {
	ip := c.RealIP()
	c.String(http.StatusOK, ip)
	return nil
}
