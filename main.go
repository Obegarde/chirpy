package main

import (
	"net/http"
	"log"	
)

func main(){
	//Create a mux
	mux := http.NewServeMux()
	//Create api config struct
	apiCfg := NewConfig() 

	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz",healthHandler)
	mux.HandleFunc("/metrics",apiCfg.hitHandler)
	mux.HandleFunc("/reset",apiCfg.resetHitHandler)

	//Create a ServerStruct
		server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Launch the server
	log.Fatal(server.ListenAndServe())
}



