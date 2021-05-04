package viciwrapper
import (
	"log"
	"context"
	"time"
	"fmt"
	"github.com/jlti-dev/ipsec_mgmt/filewrapper"
)
func (v *ViciWrapper) monitorConns(){
	lastEventTime := make(map [string]time.Time)

	v.startCommand()
	if err := v.session.Subscribe("child-updown"); err != nil {
		v.endCommand(err)
		log.Panicln(err)
		return
	}
	v.endCommand(nil)
	for {
		e, err := v.session.NextEvent(context.Background())
		if err != nil {
			log.Println(err)
			log.Panicln("Assuming vici went down, shutting down this application")
			break
		}
		k := e.Message.Keys()
		if k == nil {
			continue
		}
		log.Printf("[%s] %s\n", e.Name, k)
		for _,value := range k {
			if(value == "up"){
				//ist ein internes fragment
				continue
			}
			if value == "" || value == "(unnamed)" {
				//ignoring unnamed SAs
				continue
			}
			lastEvent, ok := lastEventTime[value]
			if (ok && time.Since(lastEvent) > 20 * time.Second){
				//v.checkChannel <- v.ikesInSystem[value]
			}
			v.checkDelay = append(v.checkDelay, value)

			lastEventTime[value] = time.Now()
		}

	}
}
func (v *ViciWrapper) watchIkes() {
	go v.monitorConns()
	log.Printf("[watch] Start watching for %d ikes\n", len(v.ikesInSystem))
	ticker := time.NewTicker(20 * time.Second)
	tickCount := 20
	for {
		select {
		case conn := <- v.terminateChannel:
			log.Printf("[%s] received to terminate\n", conn.Name)
			if errTerminate := conn.terminate(v); errTerminate != nil {
				log.Printf("[%s] could not terminate Connection: %s\n", conn.Name, errTerminate)
				continue
			}
			v.initiateChannel <- conn
		case conn := <- v.initiateChannel:
			log.Printf("[%s] received to initiate\n", conn.Name)
			if errInitiate := conn.initiateConnection(v); errInitiate != nil {
				log.Printf("[%s] could not connect: %s\n", conn.Name, errInitiate)
				continue
			}
			v.checkChannel <- v.ikesInSystem[conn.Name]
		case ike := <- v.checkChannel:
			v.checkIke(ike)
		case <- ticker.C:
			tickCount --
			for _,ike := range v.ikesInSystem {
				delayed := false
				for _, ikeDelay := range v.checkDelay{
					if ikeDelay == ike.ikeName {
						log.Printf("[WATCH] %s is in Delay Mode", ikeDelay)
						delayed = true
						break
					}
				}
				if delayed {
					//Dont check this entity now!
					continue
				}
				v.checkChannel <- ike
			}
			if len(v.checkDelay) > 0 {
				log.Printf("[WATCH] Deleting %d entries from delaylist", len(v.checkDelay))
				v.checkDelay = nil
			}

			if tickCount < 1 {
				log.Println("[WATCH] I am alive")
				tickCount = 20
			}
		}
	}
}
func (v *ViciWrapper) checkIke(ikeExpected ikeInSystem) {
	conn, errConn := v.connectionFromFile(ikeExpected.ikeName)
	if errConn != nil{
		v.UnloadConnection(ikeExpected.ikeName)
		v.ReadConnection(ikeExpected.ikeName)
		return
	}
	ike, ikeCount, err := v.findIke(ikeExpected.ikeName)
	if ikeCount == 0 && err != nil && !ikeExpected.initiator {
		log.Println(err)
		v.initiateChannel <- conn
	}else if err != nil && !ikeExpected.initiator{
		log.Println(err)
		v.terminateChannel <- conn
		return
	}
	if ikeExpected.numberRemoteTS != ike.numberRemoteTS {
		log.Printf("[WATCH][%s] Remote Traffic Selectors: expected %d, found %d\n", ikeExpected.ikeName, ikeExpected.numberRemoteTS, ike.numberRemoteTS)
	}else if ikeExpected.numberLocalTS != ike.numberLocalTS {
		log.Printf("[WATCH][%s] Local Traffic Selectors: expected %d, found %d\n", ikeExpected.ikeName, ikeExpected.numberLocalTS, ike.numberLocalTS)
	}else if ikeExpected.numberChildren != ike.numberChildren {
		log.Printf("[WATCH][%s] Children: expected %d, found %d\n", ikeExpected.ikeName, ikeExpected.numberChildren, ike.numberChildren)
	}else{
		log.Printf("[WATCH][%s] looks good!\n", ikeExpected.ikeName)
		return
	}
	if ikeExpected.initiator == false {
		log.Printf("[WATCH][%s] is not an initiator, therefor leaving check now\n", ikeExpected.ikeName)
		return
	}
	if ikeExpected.numberRemoteTS > ike.numberRemoteTS ||
	ikeExpected.numberLocalTS > ike.numberLocalTS ||
	ikeExpected.numberChildren > ike.numberChildren {
		v.initiateChannel <- conn
	}else{
		v.terminateChannel <- conn
	}
}
func (v *ViciWrapper) findIke(ikeName string)(ikeInSystem, int, error){
	retVal := ikeInSystem{
		ikeName: ikeName,
		initiator: filewrapper.GetBoolValueFromPath(ikeName, "Initiator"),
		numberRemoteTS: 0,
		numberLocalTS: 0,
		numberChildren: 0,
	}
	ikes, err := v.listSAs()
	if err != nil {
		log.Fatalf("[%s] %s", ikeName, err)
		return retVal, 0, err
	}
	ikeCnt := 0
	for _, ike := range ikes {
		if ike.Name == ikeName {
			ikeCnt ++
			retVal.Version = ike.Version
		}else {
			continue
		}
		for _, child := range ike.Children {
			retVal.numberChildren += 1
			retVal.numberRemoteTS += len(child.RemoteTS)
			retVal.numberLocalTS += len(child.LocalTS)

			selector := tsFound{
				localTS: child.LocalTS,
				remoteTS: child.RemoteTS,
			}
			retVal.selectors = append(retVal.selectors, selector)
		}
	}
	log.Printf("[CHECK] %s: %d, children: %d, remote/local: %d/%d\n", retVal.ikeName, ikeCnt, retVal.numberChildren, retVal.numberRemoteTS, retVal.numberLocalTS)
	if ikeCnt != 1 {
		return retVal,ikeCnt, fmt.Errorf("[%s] there are %d ikes connected, 1 expected!", ikeName, ikeCnt)
	}
	return retVal, ikeCnt, nil

}
