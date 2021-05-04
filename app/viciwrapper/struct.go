package viciwrapper
import (
	"time"
	"github.com/strongswan/govici/vici"
)
type ViciWrapper struct {
	session			*vici.Session
	counterCommands		int64
	counterErrors		int64
	lastCommand		time.Time
	execDuraLast		time.Duration
	execDuraAvgMs		int64
	counterSecrets		int64
	ikesInSystem		map[string]ikeInSystem
	saNameSuffix		string
	checkChannel		chan ikeInSystem
	initiateChannel		chan loadConnection
	terminateChannel	chan loadConnection
	checkDelay		[]string
}
type ViciMetrics struct {
	CounterCommands int64
	CounterErrors	int64
	LastCommand	time.Time
	ExecDuraLast	time.Duration
	ExecDuraAvgNs	int64
	LoadedSecrets	int64
}
func (v *ViciWrapper) startCommand(){
	v.lastCommand = time.Now()
}
func (v *ViciWrapper) endCommand(hasError error ){
	v.execDuraLast = time.Since(v.lastCommand)
	if hasError != nil {
		v.counterErrors ++
	}
	v.execDuraAvgMs = ( v.execDuraAvgMs * v.counterCommands + v.execDuraLast.Nanoseconds() ) / ( v.counterCommands + 1)
	v.counterCommands ++
}
type ikeInSystem struct{
	ikeName		string
	Version		int
	initiator	bool
	multiChild	bool
	numberRemoteTS	int
	numberLocalTS	int
	numberChildren	int
	selectors	[]tsFound
	duplicateTS	[]tsFound
}
type tsFound struct{
	localTS 	[]string
	remoteTS	[]string
}
type sharedSecret struct{
	Id		string			`vici:"id"`
	Typ		string			`vici:"type"`
	Data		string			`vici:"data"`
	Owners		[]string		`vici:"owners"`
}
type loadConnection struct{
	Name		string
	LocalAddrs	[]string		`vici:"local_addrs"`
	RemoteAddrs	[]string		`vici:"remote_addrs"`
	Local		AuthOpts		`vici:"local"`
	Remote		AuthOpts		`vici:"remote"`
	Children	map[string]ChildSA	`vici:"children"`
	Version		int			`vici:"version"`
	Proposals	[]string		`vici:"proposals"`
	DpdDelay	string			`vici:"dpd_delay"`
	DpdTimeout	string			`vici:"dpd_timeout"`
	Mobike		string			`vici:"mobike"`
	Encap		string			`vici:"encap"`
	RekeyTime	string			`vici:"rekey_time"`
}
type AuthOpts struct{
	Auth		string			`vici:"auth"`
	ID		string			`vici:"id"`
}
type ChildSA struct {
	Name		string
	LocalTS		[]string		`vici:"local_ts"`
	RemoteTS	[]string		`vici:"remote_ts"`
	Proposals	[]string		`vici:"esp_proposals"`
	RekeyTime	string			`vici:"rekey_time"`
}
type LoadedIKE struct {
	Name		string
	UniqueId	string			`vici:"uniqueid"`
	Version		int			`vici:"version"`
	State		string			`vici:"state"`
	LocalHost	string			`vici:"local-host"`
	RemoteHost	string			`vici:"remote-host"`
	Initiator	string			`vici:"initiator"`
	NatRemote	string			`vici:"nat-remote"`
	NatFake		string			`vici:"nat-fake"`
	EncAlg		string			`vici:"encr-alg"`
	EncKey		int			`vici:"encr-keysize"`
	IntegAlg	string			`vici:"integ-alg"`
	IntegKey	int			`vici:"integ-keysize"`
	DHGroup		string			`vici:"dh-group"`
	EstablishSec	int64			`vici:"established"`
	RekeySec	int64			`vici:"rekey-time"`
	ReauthSec	int64			`vici:"reauth-time"`
	Children	map[string]LoadedChild	`vici:"child-sas"`
}
type LoadedChild struct {
	Name		string			`vici:"name"`
	UniqueId	string			`vici:"uniqueid"`
	State		string			`vici:"state"`
	Mode		string			`vici:"mode"`
	Protocol	string			`vici:"protocol"`
	Encap		string			`vici:"encap"`
	EncAlg		string			`vici:"encr-alg"`
	EncKey		int			`vici:"encr-keysize"`
	IntegAlg	string			`vici:"integ-alg"`
	IntegKey	int			`vici:"integ-keysize"`
	DHGroup		string			`vici:"dh-group"`
	BytesIn		int64			`vici:"bytes-in"`
	PacketsIn	int64			`vici:"bytes-out"`
	LastInSec	int64			`vici:"use-in"`
	BytesOut	int64			`vici:"bytes-out"`
	PacketsOut	int64			`vici:"bytes-out"`
	LastOutSec	int64			`vici:"use-out"`
	EstablishSec	int64			`vici:"install-time"`
	RekeySec	int64			`vici:"rekey-time"`
	LifetimeSec	int64			`vici:"life-time"`
	LocalTS		[]string		`vici:"local-ts"`
	RemoteTS	[]string		`vici:"remote-ts"`
}
