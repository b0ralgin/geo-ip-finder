package controllers

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/b0ralgin/geo-ip-finder/services"
	"github.com/labstack/echo"
)

type GeoIpController struct {
	db           cache.Cache
	cacheTimeout time.Duration
	geoIpService services.Requester
}

func NewGeoController(db cache.Cache, service services.Requester, timeout time.Duration) *GeoIpController {
	return &GeoIpController{
		db:           db,
		geoIpService: service,
		cacheTimeout: timeout,
	}
}

func (g *GeoIpController) GetCountryByIp(c echo.Context) error {
	ip := c.RealIP()
	if qip := c.Param("ip"); qip != "" {
		ip = qip
	}
	if g.db.IsExist(ip) {
		return c.String(http.StatusOK, g.db.Get(ip).(string))
	}
	res, err := g.geoIpService.Request(ip)
	if err != nil {
		return c.String(http.StatusBadGateway, err.Error())
	}
	if err := g.db.Put(ip, res, g.cacheTimeout); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, res)
}
