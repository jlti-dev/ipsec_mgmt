package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"log"
	"fmt"
)

func RunWebApi(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/api/connections", getFiles).Methods("GET")
	router.HandleFunc("/api/{path}/connection", getConnection).Methods("GET")
	router.HandleFunc("/api/{path}/connection", changeSecret).Methods("UPDATE")
	router.HandleFunc("/api/{path}/secret", getSecret).Methods("GET")
	router.HandleFunc("/api/{path}/secret", changeSecret).Methods("UPDATE")
	
	router.HandleFunc("/api/{path}", createFile).Methods("POST")
	router.HandleFunc("/api/{path}", deleteFile).Methods("DELETE")

	router.HandleFunc("/api/{path}/load", loadConnection).Methods("PUT")
	router.HandleFunc("/api/{path}/unload", unloadConnection).Methods("PUT")


	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port) , router))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[webapi-request] %s: %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func getPath(r *http.Request) string {
	params := mux.Vars(r)
	return params["path"]
}
func loadConnection(w http.ResponseWriter, r *http.Request){
	load, err := ReadLoadFromFile(getPath(r))
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	load.Connection.Load()
}
func unloadConnection(w http.ResponseWriter, r *http.Request){
	load, err := ReadLoadFromFile(getPath(r))
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	load.Connection.Unload()
}
func deleteFile(w http.ResponseWriter, r *http.Request){
	path := getPath(r)
	load, err := ReadLoadFromFile(path)
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = load.Connection.Unload()
	if (err != nil) {
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = DeleteFile(path)
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func createFile(w http.ResponseWriter, r * http.Request){
	path := getPath(r)

	load := Load{}
	json.NewDecoder(r.Body).Decode(&load)

	if (load.Connection.Version == "") {
		log.Printf("[webapi] Connection not correctly transmitted")
		http.Error(w, "Connection not correct", http.StatusNotFound)
		return
	}
	err := load.Connection.Load()
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = WriteLoadToFile(&load)
	if (err != nil){
		log.Printf("[webapi] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
