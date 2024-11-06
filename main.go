package main

import (
	"net/http"
	"log"	
	_"github.com/lib/pq"		
	"github.com/joho/godotenv"
	"os"
	"database/sql"
)

func main(){


	//Load env file
	godotenv.Load()
	//get the dbUrl
	dbURL := os.Getenv("DB_URL")
	// Open a db connection
	db,err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Printf("DB error: %v",err)
	}
	defer db.Close() 

	//Create a mux
	mux := http.NewServeMux()
	//Create api config struct
	apiCfg := NewConfig(db) 

	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz",healthHandler)
	mux.HandleFunc("GET /admin/metrics",apiCfg.hitHandler)
	mux.HandleFunc("POST /admin/reset",apiCfg.resetHitHandler)
	mux.HandleFunc("POST /api/validate_chirp",handlerChirpsValidate)
	//Create a ServerStruct
		server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Launch the server
	log.Fatal(server.ListenAndServe())
}



