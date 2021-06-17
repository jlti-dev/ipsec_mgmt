package main
import (
	"log"
	"github.com/jlti-dev/ipsec_mgmt/vici"
	"os"
)
func main() {
	log.Println("Starting Application")
	
	loadLocalSecret()

	for _, path := range vici.GetAllFiles() {
		load, err := vici.ReadLoadFromFile(path)
		if (err != nil) {
			log.Println(err)
			continue
		}
		err = vici.WriteLoadToFile(load)
		if (err != nil) {
			log.Println(err)
			continue
		}
		load.Secret.Load()
		load.Connection.Load()
	}
	log.Println(vici.LoadedConnections())
}
func loadLocalSecret(){
	s := &vici.Secret{}
	s.Id = os.Getenv("LOCAL_IP")
	s.Typ = "IKE"
	s.Data = os.Getenv("LOCAL_PSK")
	s.Owners = append(s.Owners, os.Getenv("LOCAL_IP"))

	if (s.Id == "") {
		log.Fatalln("[SETUP] environment LOCAL_IP is not set")
	}else if (s.Data == ""){
		log.Fatalln("[SETUP] environment LOCAL_PSK is not set")
	}
	
	err := s.Load()
	if (err != nil) {
		log.Fatalf("[SETUP] %s", err)
	}
}
