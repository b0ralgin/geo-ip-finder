package services

import (
	"net/http"
	"time"
)

type Service struct {
	limiter chan bool
	client  *http.Client
}

type IService interface {
	DoGetRequest(url string, decoder Decoder) error
	CanAcceptRequest() bool
}

func InitService(limiter chan bool) *Service {
	service := &Service{
		limiter: limiter,
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: &http.Transport{},
		},
	}
	go service.runLimiter()
	return service
}

func (s *Service) runLimiter() {
	for _ = range s.limiter {
		time.Sleep(1 * time.Minute)
	}
}

func (s *Service) DoGetRequest(url string, decoder Decoder) error {
	resp, err := s.client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	return decoder.Decode(resp.Body)
}

func (s *Service) CanAcceptRequest() bool {
	select {
	case s.limiter <- true:
		return true
	default:
		return false
	}
}
