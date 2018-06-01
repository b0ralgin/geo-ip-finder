package services

import (
	"encoding/json"
	"fmt"
	"io"
)

type Nekudo struct {
	*Service
	url string
}

type nekudoResponse struct {
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
}

func (n *nekudoResponse) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(n)
}

func InitNekudoService(url string, limit uint16) *Nekudo {
	nekudo := &Nekudo{
		Service: InitService(make(chan bool, limit)),
		url:     url,
	}
	return nekudo
}

func (n *Nekudo) GetCountryByIp(ip string) (string, error) {
	url := fmt.Sprintf(n.url, ip)

	res := nekudoResponse{}
	err := n.DoGetRequest(url, &res)
	return res.Country.Name, err
}
