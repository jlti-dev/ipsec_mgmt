package vici
import(
	"fmt"
	"log"
	"github.com/strongswan/govici/vici"
)
func (s *Secret) Load() (error){
	v := acquireViciSession()
	defer closeViciSession(v)

	m, err := vici.MarshalMessage(s)
	if err != nil {
		return fmt.Errorf("[%s] %s\n",s.Id, err)
	}
	_, err = v.CommandRequest("load-shared", m)
	if err != nil {
		return fmt.Errorf("[%s] %s\n", s.Id, err)
	}
	log.Printf("[secret] %s was loaded\n", s.Id)
	return nil
}
func (s *Secret) Unload() (error){
	v := acquireViciSession()
	defer closeViciSession(v)

	m := vici.NewMessage()
	if err := m.Set("id", s.Id); err != nil {
		return fmt.Errorf("[unload-shared] %s\n", err)
	}
	_, err := v.CommandRequest("unload-shared", m)
	if err != nil {
		return fmt.Errorf("[unload-shared] %s\n", err)
	}
	log.Printf("[secret] %s was unloaded\n", s.Id)
	return nil
}
