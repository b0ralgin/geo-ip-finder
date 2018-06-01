package services

import (
	"errors"

	"github.com/b0ralgin/geo-ip-finder/config"
)

type GeoIpService struct {
	services []GeoIpGetter
}

type GeoIpGetter interface {
	CanAcceptRequest() bool
	GetCountryByIp(ip string) (string, error)
}

func NewGeoIpService(services map[string]config.GeoServicesCfg) *GeoIpService {
	ipstackCfg := services["ipstack"]
	ipstack := InitIpstack(ipstackCfg.URL, ipstackCfg.Token, ipstackCfg.Limit)
	return &GeoIpService{
		services: []GeoIpGetter{ipstack},
	}
}

func (s *GeoIpService) MakeRequest(ip string) (string, error) {
	for _, service := range s.services {
		if service.CanAcceptRequest() {
			return service.GetCountryByIp(ip)
		}
	}
	return "", errors.New("cannot find free Geo Ip service")
}
