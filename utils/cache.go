package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Cache interface {
	GetProxy(name string) (Proxy, error)
	Delete(name string)
	Update(name string, ip string)
}

type CacheImpl struct {
	mutex     sync.Mutex
	proxyMap  map[string]Proxy
	db        *sqlx.DB
	tableName string
}

func (c *CacheImpl) GetProxy(name string) (Proxy, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	proxy, ok := c.proxyMap[name]
	currentTime := time.Now()
	start := currentTime.Add(time.Duration(-5) * time.Minute)
	outdated := isOutDated(start, proxy.Timestamp)
	if ok && !outdated {
		fmt.Println("The cache is used")
		return proxy, nil
	}
	db := c.db
	table := c.tableName
	c.Delete(name)

	results, err := db.Query("SELECT source, endpoint, port, proxyType, user, password, apiEndpoint FROM "+table+" WHERE source = ?", name)
	if err != nil {
		fmt.Println("The failed name is: " + name)
		fmt.Println("Getting from database failed - cache")
	}
	if results.Next() {
		err = results.Scan(&proxy.ProxyName, &proxy.Endpoint, &proxy.Port, &proxy.ProxyType, &proxy.User, &proxy.Pass, &proxy.APIEndpoint)
		if err != nil {
			fmt.Println("Scanning from rows failed")
		}
		currentTime := time.Now()
		proxy.Timestamp = currentTime
		c.Update(proxy.ProxyName, proxy)
		fmt.Println("The cache is updated")
	}

	if proxy.ProxyType == "dynamic" {
		ips, err := GetIPFromAPI(proxy.APIEndpoint)
		ip := ips[(rand.Intn(10))]
		if err != nil {
			fmt.Println("Getting IP from external api failed")
		}
		host, port, err := net.SplitHostPort(ip)
		if err != nil {
			fmt.Println("Splitting host and port failed")
		}
		proxy.Endpoint = host
		portNum, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(err)
		}
		proxy.Port = portNum
		currentTime = time.Now()
		proxy.Timestamp = currentTime
		c.Update(proxy.ProxyName, proxy)
	}
	return proxy, nil
}

func (c *CacheImpl) Delete(name string) {
	delete(c.proxyMap, name)
}

func (c *CacheImpl) Update(name string, proxy Proxy) {
	c.proxyMap[name] = proxy
}

func NewCache(sqlConn string, tableName string) (*CacheImpl, error) {
	// assume that tableName is valid
	fmt.Println(sqlConn)
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", sqlConn)
	if db == nil {
		return nil, err
	}

	return &CacheImpl{
		db:        db,
		tableName: tableName,
		mutex:     sync.Mutex{},
		proxyMap:  map[string]Proxy{},
	}, nil
}

func isOutDated(start, check time.Time) bool {
	return check.Before(start)
}
