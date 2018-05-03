package services

import (
	"fmt"
	logger "log"
	"runtime/debug"
)

type Service interface {
	Start(args interface{}, log *logger.Logger) (err error)
	Clean()
}
type ServiceItem struct {
	S    Service
	Args interface{}
	Name string
	Log  *logger.Logger
}

var servicesMap = map[string]*ServiceItem{}

func Regist(name string, s Service, args interface{}, log *logger.Logger) {
	Stop(name)
	servicesMap[name] = &ServiceItem{
		S:    s,
		Args: args,
		Name: name,
		Log:  log,
	}
}
func GetService(name string) *ServiceItem {
	if s, ok := servicesMap[name]; ok && s.S != nil {
		return s
	}
	return nil

}
func Stop(name string) {
	if s, ok := servicesMap[name]; ok && s.S != nil {
		s.S.Clean()
	}
}
func Run(name string, args interface{}) (service *ServiceItem, err error) {
	service, ok := servicesMap[name]
	if ok {
		defer func() {
			e := recover()
			if e != nil {
				err = fmt.Errorf("%s servcie crashed, ERR: %s\ntrace:%s", name, e, string(debug.Stack()))
			}
		}()
		if args != nil {
			err = service.S.Start(args, service.Log)
		} else {
			err = service.S.Start(service.Args, service.Log)
		}
		if err != nil {
			err = fmt.Errorf("%s servcie fail, ERR: %s", name, err)
		}
	} else {
		err = fmt.Errorf("service %s not found", name)
	}
	return
}
