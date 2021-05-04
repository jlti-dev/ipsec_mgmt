package filewrapper

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
	"io/ioutil"
	"log"
)
func getAllFiles() []string{
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
func GetFilesForSecrets() []string{
	var files []string
	for _, file := range getAllFiles() {
		if strings.HasSuffix(file, ".secret"){
			files = append(files, file)
		}
	}
	return files
}
func GetFilesForConnections() []string{
	var files []string
	for _, file := range getAllFiles() {
		if strings.HasSuffix(file, ".secret"){
			continue
		}
		files = append(files, file)
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
func WriteOrReplaceLine(path string, prefix string, value string) error {
	file := strings.Join([]string{"/app", "config", path, }, "/")
	input, err := ioutil.ReadFile(file)
	prefixWritten := false
	var lines []string
	if err != nil {
		log.Printf("[file] %s does not exist, creating it", path)
	}else{
		lines = strings.Split(string(input), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, prefix) {
				lines[i] = strings.Join([]string{prefix, value,}, "=")
				prefixWritten = true
				break
			}
		}
	}
	if prefixWritten == false {
		lines = append(lines, strings.Join([]string{prefix, value, }, "="))
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0666)
	if err != nil {
		return fmt.Errorf("[%s - file] %s", path, err)
	}
	log.Printf("[%s - file] prefix %s was written successfully\n", path, prefix)
	return nil
}
func GetStringValueFromPath(path string, value string) string{
	f, err := os.Open(strings.Join([]string{"/app","config",path}, "/"))
	if err != nil{
		fmt.Println(err)
		return ""
	}
	defer func(){
		if err = f.Close(); err != nil{
			fmt.Println(err)
		}
	}()
	s := bufio.NewScanner(f)
	for s.Scan(){
		if(strings.HasPrefix(s.Text(), value)){
			return strings.Split(s.Text(), "=")[1]
		}
	}
	return "";
}
func GetBoolValueFromPath(path string, value string) bool {
	v := GetStringValueFromPath(path,value)
	if v == "yes" || v == "true" {
		return true
	}
	return false
}
func GetIntValueFromPath(path string, value string) int {
	i, err := strconv.Atoi(GetStringValueFromPath(path, value))
	if err != nil {
		fmt.Printf("Error in file %s, value %s\n", path, value)
		fmt.Println(err)
		return 0
	}
	return i
}
func GetStringArrayFromPath(path string, value string) []string {
	return strings.Split(GetStringValueFromPath(path, value),",")
}
