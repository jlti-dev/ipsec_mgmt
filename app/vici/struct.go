package vici
type Load struct {
	Connection	Connection	`json:"connections,omitempty"`
	Secret		Secret		`json:"secret,omitempty"`
	Path		string			`json:"-"`
}
type Connection struct {
	Name		string			`json:"name"`
	LocalAddr	[]string		`json:"local_addrs,omitempty" vici:"local_addrs"`
	Remote_addrs	[]string		`json:"remote_addrs,omitempty" vici:"remote_addrs"`
	Local		Auth			`json:"local,omitempty" vici:"local"`
	Remote		Auth			`json:"remote,omitempty" vici:"remote"`
	Children	map[string]ChildSA	`json:"children,omitempty" vici:"children"`
	Version		string			`json:"version,omitempty" vici:"version"`
	Proposals	[]string		`json:"proposals,omitempty" vici:"proposals"`
	DpdDelay	string			`json:"dpd_delay,omitempty" vici:"dpd_delay"`
	DpdTimeout	string			`json:"dpd_timeout,omitempty" vici:"dpd_timeout"`
	Mobike		string			`json:"mobike,omitempty" vici:"mobike"`
	Encap		string			`json:"encap,omitempty" vici:"encap"`
	RekeyTime	string			`json:"rekey_time,omitempty" vici:"rekey_time"`
}
type Auth struct {
	Auth		string			`json:"auth,omitempty" vici:"auth"`
	ID		string			`json:"id,omitempty" vici:"id"`
}
type ChildSA struct {
	LocalTS		[]string		`json:"local_ts,omitempty" vici:"local_ts"`
	RemoteTS	[]string		`json:"remote_ts,omitempty" vici:"remote_ts"`
	Proposals	[]string		`json:"esp_proposals,omitempty" vici:"esp_proposals"`
	RekeyTime	string			`json:"rekey_time,omitempty" vici:"rekey_time"`
	ReadLocalTS	[]string		`vici:"local-ts"`
	ReadRemoteTS	[]string		`vici:"remote-ts"`
	StartAction	string			`json:"start_action" vici:"start_action"`
	CloseAction	string			`json:"close_action" vici:"close_action"`
}
type Secret struct{
	Id		string			`json:"id,omitempty" vici:"id"`
	Typ		string			`json:"type,omitempty" vici:"type"`
	Data		string			`json:"data,omitempty" vici:"data"`
	Owners		[]string		`json:"owners,omitempty" vici:"owners"`
}
