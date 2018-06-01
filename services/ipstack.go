package services

import (
	"log"
	"net/url"
	"time"
)

type Ipstack struct {
	api, token string
	limiter    chan bool
}

func InitIpstack(api, token string, limit uint16) *Ipstack {
	ipstack := &Ipstack{
		api:     api,
		token:   token,
		limiter: make(chan bool, limit),
	}
	go ipstack.setupLimiter()
	return ipstack
}

func (i *Ipstack) setupLimiter() {
	for {
		select {
		case <-i.limiter:
			time.Sleep(1 * time.Minute)
		}
	}
}

func (i *Ipstack) getReqestUrl(ip string) (string, error) {
	u, err := url.Parse("/" + ip)
	if err != nil {
		log.Fatal(err)
	}

	queryString := u.Query()
	queryString.Set("token", i.token)
	u.RawQuery = queryString.Encode()

	apiURL, err := url.Parse(i.api)
	if err != nil {
		return "", nil
	}

	return apiURL.ResolveReference(u).String(), nil
}

func (i *Ipstack) GetCountryByIp(ip string) (string, error) {
	log.Println(i.getReqestUrl(ip))
	return ip, nil
}

func (i *Ipstack) CanAcceptRequest() bool {
	select {
	case i.limiter <- true:
		return true
	default:
		return false
	}
}
