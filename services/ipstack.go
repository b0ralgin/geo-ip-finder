package services

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
)

type Ipstack struct {
	*Service
	api, token string
}

type ipstackResponse struct {
	CountryName string `json:"country_name"`
}

func (i *ipstackResponse) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(i)
}

func InitIpstack(api, token string, limit uint16) *Ipstack {
	return &Ipstack{
		Service: InitService(make(chan bool, limit)),
		api:     api,
		token:   token,
	}
}

func (i *Ipstack) getReqestUrl(ip string) (string, error) {
	u, err := url.Parse("/" + ip)
	if err != nil {
		log.Fatal(err)
	}

	queryString := u.Query()
	queryString.Set("access_key", i.token)
	u.RawQuery = queryString.Encode()

	apiURL, err := url.Parse(i.api)
	if err != nil {
		return "", nil
	}

	return apiURL.ResolveReference(u).String(), nil
}

func (i *Ipstack) GetCountryByIp(ip string) (string, error) {
	url, err := i.getReqestUrl(ip)
	if err != nil {
		return "", err
	}
	res := ipstackResponse{}
	err = i.DoGetRequest(url, &res)
	return res.CountryName, err
}
