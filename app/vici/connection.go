package vici

import(
	"log"
	"github.com/strongswan/govici/vici"
	"fmt"
)

func (c *Connection) Load() (error){
	log.Printf("[Connection.Load] requested to load connection name %s\n", c.Name)
	v := acquireViciSession()
	defer closeViciSession(v)

	msg, err := vici.MarshalMessage(c)
	if (err != nil){
		return fmt.Errorf("[Connection.Load] %s", err)
	}
	m := vici.NewMessage()
	m.Set(c.Name, msg)
	
	_, err = v.CommandRequest("load-conn", m)
	if (err != nil){
		return fmt.Errorf("[Connection.Load] %s", err)
	}


	log.Printf("[Connection.Load] loaded connection as %s\n", c.Name)
	return nil
}
func (c *Connection) Unload() (error){
	log.Printf("[Conenction.Unload] requested to unload connection name %s\n", c.Name)
	v := acquireViciSession()
	defer closeViciSession(v)

	m := vici.NewMessage()
	
	err := m.Set("name", c.Name)
	if (err != nil){
		return fmt.Errorf("[Connection.Unload] %s", err)
	}

	
	_, err = v.CommandRequest("load-conn", m)
	if (err != nil){
		return fmt.Errorf("[Connection.Unload] %s", err)
	}


	log.Printf("[Connection.Unload] unloaded connection as %s\n", c.Name)
	return nil
}
func LoadedConnections()(map[string]Connection){
	retVar := map[string]Connection{}

	v := acquireViciSession()
	defer closeViciSession(v)

	msgs, err := v.StreamedCommandRequest("list-conns", "list-conn", nil)
	if err != nil {
		log.Printf("[LoadedConnections] %s", err)
		return retVar
	}
	for _,m := range msgs.Messages() {
		if e := m.Err(); e != nil{
			//ignoring this error
			continue
		}
		for _, k := range m.Keys() {
			inbound := m.Get(k).(*vici.Message)
			var conn Connection
			log.Println(inbound)
			if e := vici.UnmarshalMessage(inbound, &conn); e != nil {
				//ignoring this marshal/unmarshal errro!
				log.Println(e)
				continue
			}
			conn.Name = k
			retVar[k] = conn
		}
	}
	return retVar
}

