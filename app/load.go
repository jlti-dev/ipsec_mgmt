package main

import(
	"fmt"
	"log"
	"encoding/json"
)

func ReadLoadFromFile(path string) (*Load, error){
	log.Printf("[ReadLoadFromFile] accessing %s\n", path)
	in, err := ReadFile(path)
	if (err != nil){
		return nil, fmt.Errorf("[ReadLoadFromFile] %s", err)
	}
	l := &Load{}
	err = json.Unmarshal([]byte(in), &l)
	if (err != nil){
		return nil, fmt.Errorf("[ReadLoadFromFile] %s", err)
	}
	l.Path = path
	log.Printf("[ReadLoadFromFile] %s was read\n", path)
	return l, nil
}
func WriteLoadToFile(load *Load)(error){
	log.Printf("[WriteLoadToFile] accessing %s\n", load.Path)
	l, err := json.Marshal(load)
	if err != nil {
		return fmt.Errorf("[WriteLoadToFile] %s", err)
	}
	err = WriteFile(load.Path, string(l))
	if err != nil {
		return fmt.Errorf("[WriteLoadToFile] %s", err)
	}
	log.Printf("[WriteLoadToFile] %s was written\n", load.Path)
	return nil
}
