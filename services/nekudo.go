package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Nekudo struct {
	url     string
	client  *http.Client
	limiter chan bool
}

type nekudoResponse struct {
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
}

func InitNekudoService(url string, limit uint16) *Nekudo {
	nekudo := &Nekudo{
		url:     url,
		limiter: make(chan bool, limit),
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: &http.Transport{},
		},
	}
	go nekudo.runLimiter()
	return nekudo
}

func (n *Nekudo) runLimiter() {
	for {
		select {
		case <-n.limiter:
			time.Sleep(1 * time.Minute)
		}
	}
}

func (n *Nekudo) CanAcceptRequest() bool {
	select {
	case n.limiter <- true:
		return true
	default:
		return false
	}
}

func (n *Nekudo) GetCountryByIp(ip string) (string, error) {
	url := fmt.Sprintf(n.url, ip)
	resp, err := n.client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	res := nekudoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	return res.Country.Name, err
}
