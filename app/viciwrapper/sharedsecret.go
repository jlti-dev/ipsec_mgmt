package viciwrapper

import (
	"fmt"
	"log"
	"github.com/jlti-dev/ipsec_mgmt/filewrapper"
	"github.com/strongswan/govici/vici"
)
func (v *ViciWrapper) countSecrets() (int, error) {
	v.startCommand()
	loaded, e := v.session.CommandRequest("get-shared", nil)
	v.endCommand(e)
	if e != nil {
		return 0, fmt.Errorf("[get-shared] %s\n", e)
	}

	return len(loaded.Get("keys").([]string)), nil
}
func (v *ViciWrapper) isSecretLoaded( secretId string) (bool, error){
	v.startCommand()
	loaded, e := v.session.CommandRequest("get-shared", nil)
	v.endCommand(e)
	if e != nil {
		return false, fmt.Errorf("[get-shared] %s\n", e)
	}
	for _, value := range loaded.Get("keys").([]string){
		if value == secretId {
			return true, nil
		}
	}
	return false, nil
}
func (v *ViciWrapper) unloadSecret(secretId string) error{
	m := vici.NewMessage()
	if err := m.Set("id", secretId); err != nil {
		return fmt.Errorf("[unload-shared] %s\n", err)
	}
	v.startCommand()
	_, e := v.session.CommandRequest("unload-shared", m)
	v.endCommand(e)
	if e != nil {
		return fmt.Errorf("[unload-shared] %s\n", e)
	}
	log.Printf("[secret] %s was unloaded\n", secretId)
	return nil

}
func (v *ViciWrapper) loadSharedSecret(path string) error{
	psk := sharedSecret{
		Id: filewrapper.GetStringValueFromPath(path, "RemoteAddrs"),
		Typ: "IKE",
		Data: filewrapper.GetStringValueFromPath(path, "PSK"),
		Owners: filewrapper.GetStringArrayFromPath(path, "RemoteAddrs"),
	}
	if psk.Data == "" {
		return fmt.Errorf("Secret in file %s is no PSK\n", path)
	}
	isLoaded, err := v.isSecretLoaded(psk.Id)
	if err != nil {
		return err
	}else if isLoaded {
		v.unloadSecret(psk.Id)
	}
	m, err := vici.MarshalMessage(psk)
	if err != nil {
		return fmt.Errorf("[%s] %s\n",path, err)
	}
	v.startCommand()
	_, err2 := v.session.CommandRequest("load-shared", m)
	v.endCommand(err2)
	if err2 != nil {
		return fmt.Errorf("[%s] %s\n", path, err2)
	}
	log.Printf("[secret] %s was loaded\n", path)
	return nil
}
