package services

import (
	"log"
)

type UDP struct {
}

func NewUDP() Service {
	return &UDP{}
}
func (s *UDP) Start(args interface{}) (err error) {
	log.Printf("called")
	return
}
func (s *UDP) Clean() {

}
