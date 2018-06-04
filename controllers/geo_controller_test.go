package controllers

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type service struct {
	service map[string]string
}

func (s service) Request(ip string) (string, error) {
	return s.service[ip], nil
}

type failedService struct {
	service map[string]string
}

func (s failedService) Request(ip string) (string, error) {
	return "", errors.New("Error")
}

func TestGetCountryByIpAsParam(t *testing.T) {
	ip := "1.1.1.1"
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:ip")
	c.SetParamNames("ip")
	c.SetParamValues(ip)
	s := service{
		service: map[string]string{
			ip: "Russia",
		},
	}
	ipCache := cache.NewMemoryCache()
	defer ipCache.ClearAll()
	g := NewGeoController(ipCache, s, time.Duration(1*time.Second))
	if assert.NoError(t, g.GetCountryByIp(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Russia", rec.Body.String())
	}
}

func TestGetCountryByIp(t *testing.T) {
	ip := "192.0.2.1"
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s := service{
		service: map[string]string{
			ip: "Russia",
		},
	}
	ipCache := cache.NewMemoryCache()
	defer ipCache.ClearAll()
	g := NewGeoController(ipCache, s, time.Duration(1*time.Second))
	if assert.NoError(t, g.GetCountryByIp(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Russia", rec.Body.String())
	}
}

func TestGetCountryByIpBeInCache(t *testing.T) {
	ip := "192.0.2.1"
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s := service{
		service: map[string]string{
			ip: "Russia",
		},
	}
	ipCache := cache.NewMemoryCache()
	defer ipCache.ClearAll()
	g := NewGeoController(ipCache, s, time.Duration(10*time.Second))
	if assert.NoError(t, g.GetCountryByIp(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Russia", rec.Body.String())
	}
	rec.Body.Reset()
	if assert.NoError(t, g.GetCountryByIp(c)) {
		log.Println("Resp", ipCache.Get(ip).(string))
		assert.Equal(t, ipCache.Get(ip).(string), "Russia")
		assert.Equal(t, "Russia", rec.Body.String())
	}
}

func TestGetCountryByIpGetErrorWhileRequest(t *testing.T) {
	ip := "192.0.2.1"
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s := failedService{
		service: map[string]string{
			ip: "Russia",
		},
	}
	ipCache := cache.NewMemoryCache()
	defer ipCache.ClearAll()
	g := NewGeoController(ipCache, s, time.Duration(10*time.Second))
	if assert.NoError(t, g.GetCountryByIp(c)) {
		assert.Equal(t, http.StatusBadGateway, rec.Code)
	}
}
