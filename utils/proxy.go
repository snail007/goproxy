package utils

import "time"

type Proxy struct {
	ProxyName   string
	Endpoint    string
	Port        int
	ProxyType   string
	User        string
	Pass        string
	Timestamp   time.Time
	APIEndpoint string
}

func (p *Proxy) SetIP(endpoint string, port int) {
	p.Endpoint = endpoint
	p.Port = port
}

func NewProxy(name string, endpoint string, port int, proxyType string, user string, pass string, timestamp time.Time, apiEndPoint string) *Proxy {
	return &Proxy{
		ProxyName:   name,
		Endpoint:    endpoint,
		Port:        port,
		ProxyType:   proxyType,
		User:        user,
		Pass:        pass,
		Timestamp:   timestamp,
		APIEndpoint: apiEndPoint,
	}
}
