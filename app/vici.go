package main

import(
	"github.com/strongswan/govici/vici"
	"log"
	"fmt"
)

func acquireViciSession() (*vici.Session){
	s, err := vici.NewSession()
	if (s != nil) {
		log.Println("[acquireViciSession] opening session")
		return s
	}else if (err != nil) {
		err = fmt.Errorf("[acquireViciSession] %s\n", err)
	}else {
		err = fmt.Errorf("[acquireViciSession] not possible")
	}
	if (err != nil) {
		log.Fatalf("[acquireViciSession] %s\n", err)
	}
	return s
}
func closeViciSession(s *vici.Session) () {
	if (s != nil) {
		log.Println("closing Session")
		err := s.Close()
		if (err != nil){
			log.Printf("[closeViciSession] %s\n", err)
		}
	}else{
		log.Printf("[closeViciSession] viciSession is nil, cant close")
	}
}

