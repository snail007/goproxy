package utils

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Pattern struct {
	Patt      string
	ProxyName string
	Priority  int
}

func NewPattern(pattern string, proxyName string, priority int) *Pattern {
	return &Pattern{
		Patt:      pattern,
		ProxyName: proxyName,
		Priority:  priority,
	}
}

type PatternTable interface {
	Get(url string) (string, bool)
	GetFromCache(url string) (string, bool)
	Delete(patternString string)
}

type PatternTableImpl struct {
	mutex     sync.Mutex
	db        *sqlx.DB
	tableName string
	cache     []Pattern
	timestamp time.Time
}

func NewPatternTable(sqlConn string, tableName string) (*PatternTableImpl, error) {
	// assume that tableName is valid
	fmt.Println(sqlConn)
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", sqlConn)
	if db == nil {
		return nil, err
	}

	return &PatternTableImpl{
		db:        db,
		tableName: tableName,
		mutex:     sync.Mutex{},
		cache:     []Pattern{},
		timestamp: time.Now(),
	}, nil
}

func (pt *PatternTableImpl) Get(url string) (string, bool) {
	patternFromCache, ok := pt.GetFromCache(url)
	if ok {
		return patternFromCache, true
	}

	results, err := pt.db.Query("SELECT pattern, proxyName, priority FROM " + pt.tableName)
	var tempPattern Pattern
	resultPattern := NewPattern("", "", 0)
	var max int = 0
	if err != nil {
		fmt.Println("Getting from database failed - pattern table")
	}
	for results.Next() {
		err = results.Scan(&tempPattern.Patt, &tempPattern.ProxyName, &tempPattern.Priority)
		if err != nil {
			fmt.Println("Scanning from rows failed")
		}
		match, _ := regexp.MatchString(tempPattern.Patt, url)
		if match && tempPattern.Priority > max {
			*resultPattern = tempPattern
			max = tempPattern.Priority
		}
	}
	if resultPattern.ProxyName == "" {
		return "", false
	}
	pt.cache = append(pt.cache, *resultPattern)
	return resultPattern.ProxyName, true
}

func (pt *PatternTableImpl) GetFromCache(url string) (string, bool) {
	currentTime := time.Now()
	start := currentTime.Add(time.Duration(-5) * time.Minute)
	outdated := isOutDated(start, pt.timestamp)
	if outdated {
		pt.cache = pt.cache[:0]
		pt.timestamp = time.Now()
		return "", false
	}

	var max int = 0
	var proxyName string = ""
	for _, proxy := range pt.cache {
		match, _ := regexp.MatchString(proxy.Patt, url)
		if match && proxy.Priority > max {
			proxyName = proxy.ProxyName
			max = proxy.Priority
		}
	}
	if proxyName != "" {
		return proxyName, true
	}
	return "", false
}

func (pt *PatternTableImpl) Delete(patternString string) {
	_, err := pt.db.NamedExec("DELETE FROM "+pt.tableName+" WHERE pattern = :pattern", map[string]interface{}{
		"pattern": patternString,
	})
	if err != nil {
		fmt.Println("Delete pattern failed")
	}
}
