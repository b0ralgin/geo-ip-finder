package main

import (
	"flag"
	"log"

	"github.com/astaxie/beego/cache"
	"github.com/b0ralgin/geo-ip-finder/config"
	"github.com/b0ralgin/geo-ip-finder/controllers"
	"github.com/b0ralgin/geo-ip-finder/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	cnfFilename = flag.String("config", "config/config.yaml", "Filename of config file (YML format)")
)

func main() {
	flag.Parse()

	e := echo.New()
	config, err := config.LoadCfg(*cnfFilename)
	if err != nil {
		log.Fatal("cannot parse config file:", err)
	}
	bm, err := cache.NewCache("memory", `{"interval":-1}`)
	if err != nil {
		log.Fatal(err)
	}

	geoService := services.NewGeoIpService(config.Services)
	geoController := controllers.NewGeoController(bm, geoService, config.TTL)
	e.Use(middleware.Logger())

	e.GET("/", geoController.GetCountryByIp)
	e.GET("/:ip", geoController.GetCountryByIp)

	e.Logger.Fatal(e.Start(":1323"))
}
