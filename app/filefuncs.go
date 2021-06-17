package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"log"
)
func GetAllFiles() []string{
	var files []string
	f, err := os.Open("/app/config")
	if err != nil {
		fmt.Println(err)
		return files
	}
	defer f.Close()
	fileInfo, err := f.Readdir(-1)
	if err != nil{
		fmt.Println(err)
		return files
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files
}
func DeleteFile(path string) error {
	file := strings.Join([]string{"/app", "config", path, }, "/")
	err := os.Remove(file)
	if err != nil {
		log.Printf("[file] %s could not be deleted", path)
	}
	return err
}
func ReadFile(path string) (string, error){
	file := strings.Join([]string{"/app", "config", path, }, "/")
	ret, err := ioutil.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("[%s - file] %s", path, err)
	}
	log.Printf("[%s - file] read successfully\n", path)
	return string(ret), err
}
func WriteFile(path string, value string) error {
	file := strings.Join([]string{"/app", "config", path, }, "/")
	err := ioutil.WriteFile(file, []byte(value), 0666)
	if err != nil {
		return fmt.Errorf("[%s - file] %s", path, err)
	}
	log.Printf("[%s - file] written successfully\n", path)
	return nil
}
