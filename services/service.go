package services

import (
	"fmt"
	"runtime/debug"
)

type Service interface {
	Start(args interface{}) (err error)
	Clean()
}
type ServiceItem struct {
	S    Service
	Args interface{}
	Name string
}

var servicesMap = map[string]*ServiceItem{}

func Regist(name string, s Service, args interface{}) {
	servicesMap[name] = &ServiceItem{
		S:    s,
		Args: args,
		Name: name,
	}
}
func Run(name string, args ...interface{}) (service *ServiceItem, err error) {
	service, ok := servicesMap[name]
	if ok {
		defer func() {
			e := recover()
			if e != nil {
				err = fmt.Errorf("%s servcie crashed, ERR: %s\ntrace:%s", name, e, string(debug.Stack()))
			}
		}()
		if len(args) == 1 {
			err = service.S.Start(args[0])
		} else {
			err = service.S.Start(service.Args)
		}
		if err != nil {
			err = fmt.Errorf("%s servcie fail, ERR: %s", name, err)
		}
	} else {
		err = fmt.Errorf("service %s not found", name)
	}
	return
}
