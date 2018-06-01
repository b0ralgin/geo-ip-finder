package services

import (
	"errors"
	"io"

	"github.com/b0ralgin/geo-ip-finder/config"
)

type GeoIpService struct {
	services []GeoIpGetter
}

type MakeRequester interface {
	MakeRequest(ip string) (string, error)
}

type GeoIpGetter interface {
	CanAcceptRequest() bool
	GetCountryByIp(ip string) (string, error)
}

type Decoder interface {
	Decode(io.Reader) error
}

func NewGeoIpService(services map[string]config.GeoServicesCfg) *GeoIpService {
	ipstackCfg := services["ipstack"]
	ipstack := InitIpstack(ipstackCfg.URL, ipstackCfg.Token, ipstackCfg.Limit)
	nekudoCfg := services["nekudo"]
	nekudo := InitNekudoService(nekudoCfg.URL, nekudoCfg.Limit)
	return &GeoIpService{
		services: []GeoIpGetter{
			ipstack,
			nekudo,
		},
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
