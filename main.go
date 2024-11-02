package main

import (
	"net/http"
	"log"
)

func main(){
	//Create a mux
	mux := http.NewServeMux()
	
	mux.Handle("/app/",http.StripPrefix("/app",http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz",healthHandler)
	//Create a ServerStruct
		server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Launch the server
	log.Fatal(server.ListenAndServe())
}

func healthHandler(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
