package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeService struct {
	acceptRequest bool
	db            map[string]string
}

func (f *fakeService) CanAcceptRequest() bool {
	return f.acceptRequest
}

func (f *fakeService) GetCountryByIp(ip string) (string, error) {
	return f.db[ip], nil
}
func TestMakeRequest(t *testing.T) {
	testIp := "1.1.1.1"
	service := &fakeService{
		acceptRequest: true,
		db: map[string]string{
			testIp: "Russia",
		},
	}
	geoIpService := GeoIpService{services: []GeoIpGetter{
		service,
	}}
	res, err := geoIpService.Request(testIp)
	if assert.NoError(t, err) {
		assert.Equal(t, "Russia", res)
	}
}

func TestMakeRequestError(t *testing.T) {
	testIp := "1.1.1.1"
	service := &fakeService{
		acceptRequest: false,
		db: map[string]string{
			testIp: "Russia",
		},
	}
	geoIpService := GeoIpService{services: []GeoIpGetter{
		service,
	}}
	res, err := geoIpService.Request(testIp)
	if assert.Error(t, err) {
		assert.Equal(t, "", res)
		assert.Equal(t, "cannot find free Geo Ip service", err.Error())
	}
}

func TestMakeRequest2Services(t *testing.T) {
	testIp := "1.1.1.1"
	service1 := &fakeService{
		acceptRequest: false,
		db: map[string]string{
			testIp: "Russia",
		},
	}
	service2 := &fakeService{
		acceptRequest: true,
		db: map[string]string{
			testIp: "Russian Federation",
		},
	}
	geoIpService := GeoIpService{services: []GeoIpGetter{
		service1,
		service2,
	}}
	res, err := geoIpService.Request(testIp)
	if assert.NoError(t, err) {
		assert.Equal(t, "Russian Federation", res)
	}
}
