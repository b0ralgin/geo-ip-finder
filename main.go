package main

import (
	"flag"
	"log"

	"github.com/b0ralgin/geo-ip-finder/config"
	"github.com/b0ralgin/geo-ip-finder/controllers"
	"github.com/b0ralgin/geo-ip-finder/services"
	"github.com/boltdb/bolt"
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
	db, err := bolt.Open(config.DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	geoService := services.NewGeoIpService(config.Services)
	geoController := controllers.NewGeoController(db, geoService)
	e.Use(middleware.Logger())

	e.GET("/", geoController.GetCountryByIp)

	e.Logger.Fatal(e.Start(":1323"))
}
