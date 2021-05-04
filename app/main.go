package main
import (
	"log"
	"os/signal"
	"os"
	"syscall"
	"github.com/jlti-dev/ipsec_mgmt/viciwrapper"
	"github.com/jlti-dev/ipsec_mgmt/filewrapper"
)
func main() {
	log.Println("Starting Application")
	vici, err := viciwrapper.GetWrapper()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Application initiated, reading Secrets")
	for _, secret := range filewrapper.GetFilesForSecrets() {
		e := vici.ReadSecret(secret)
		if e != nil {
			log.Printf("[%s] Shared Secret not load: %s\n", secret, e)
		}else{
			log.Printf("[%s] Shared Secret loaded\n", secret)
		}
	}
	log.Println("Reading Secrets finished")
	log.Println("Reading Connections")
	for _, conn := range filewrapper.GetFilesForConnections() {
		e := vici.ReadConnection(conn)
		if e != nil {
			log.Printf("[%s] Connection not loaded: %s\n", conn, e)
		}else{
			log.Printf("[%s] Connection loaded\n", conn)
		}
	}
	log.Println("Reading Connections finished")


	log.Println("Vici loaded, starting operations")
	osCall := make(chan os.Signal, 1)
	signal.Notify(osCall, os.Interrupt, syscall.SIGTERM)

	go vici.WatchIkes()
	for {
		select {
		case <- osCall:
			log.Fatalln("Shutdown received")
			break
		}
	}
}
