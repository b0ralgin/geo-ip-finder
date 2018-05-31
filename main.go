package main

import (
	"log"

	"github.com/b0ralgin/geo-ip-finder/controllers"
	"github.com/b0ralgin/geo-ip-finder/services"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	geoService := services.NewGeoIpService()
	geoController := controllers.NewGeoController(db, geoService)
	e.Use(middleware.Logger())

	e.GET("/", geoController.GetIp)

	e.Logger.Fatal(e.Start(":1323"))
}
