package viciwrapper
import (
	"github.com/strongswan/govici/vici"
	"github.com/jlti-dev/ipsec_mgmt/filewrapper"
	"log"
)
var me *ViciWrapper

func GetWrapper() (*ViciWrapper, error) {
	if me != nil {
		return me, nil
	}
	//Singleton not yet created:
	me = &ViciWrapper{}
	me.startCommand()
	s, err := vici.NewSession()
	me.endCommand(err)
	if err != nil {
		return &ViciWrapper{}, err
	}
	me.session = s
	me.ikesInSystem = make(map[string]ikeInSystem)
	me.checkChannel = make(chan ikeInSystem, 100)
	me.terminateChannel = make(chan loadConnection, 10)
	me.initiateChannel = make(chan loadConnection, 10)
	me.saNameSuffix = "net"
	return me, nil
}
func (w *ViciWrapper) GetViciMetrics() ViciMetrics{
	secrets, err := w.countSecrets()
	if err != nil {
		log.Println(err)
		secrets = 0
	}
	return ViciMetrics{
		CounterCommands: w.counterCommands,
		CounterErrors: w.counterErrors,
		LastCommand: w.lastCommand,
		ExecDuraLast: w.execDuraLast,
		ExecDuraAvgNs: w.execDuraAvgMs,
		LoadedSecrets: int64(secrets),
	}
}
func (w *ViciWrapper) ReadSecret(pathToFile string) error {
	return w.loadSharedSecret(pathToFile)
}
func (w *ViciWrapper) UnloadSecret(pathToFile string) error {
	return w.unloadSecret(filewrapper.GetStringValueFromPath(pathToFile, "RemoteAddrs"))
}
func (w *ViciWrapper) UnloadConnection(pathToFile string) error {
	conn, err := w.connectionFromFile(pathToFile)
	if err != nil {
		return err
	}
	ikes := []ikeInSystem{}
	for _, ike := range w.ikesInSystem {
		if ike.ikeName == pathToFile {
			continue
		}
		ikes = append(ikes, ikeInSystem{
			ikeName: pathToFile,
		})
	}
	return conn.unloadConnection(w)

}
func (w *ViciWrapper) ReadConnection(pathToFile string) error {
	_, err := w.loadConn(pathToFile)
	return err
}
func (w *ViciWrapper) ListIkes()([]LoadedIKE, error){
	return w.listSAs()
}
func (w *ViciWrapper) WatchIkes(){
	w.watchIkes()
}
func (w *ViciWrapper) GetIkesInSystem() int {
	return len(w.ikesInSystem)
}
