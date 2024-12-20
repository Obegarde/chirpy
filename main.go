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
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	// Open a db connection
	db,err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Printf("DB error: %v",err)
	}
	defer db.Close() 

	//Create a mux
	mux := http.NewServeMux()
	//Create api config struct
	apiCfg := NewConfig(db,platform,secret,polkaKey) 

	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz",healthHandler)
	mux.HandleFunc("GET /admin/metrics",apiCfg.hitHandler)
	mux.HandleFunc("POST /admin/reset",apiCfg.resetUsersHandler)
	mux.HandleFunc("POST /api/chirps",apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	mux.HandleFunc("GET /api/chirps",apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}",apiCfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}",apiCfg.handlerDeleteChirpByID)
	mux.HandleFunc("POST /api/login",apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh",apiCfg.handlerCheckRefresh)
	mux.HandleFunc("POST /api/revoke",apiCfg.handlerRevokeRefresh)
	mux.HandleFunc("POST /api/polka/webhooks",apiCfg.handlerUpgradeUser)
	//Create a ServerStruct
		server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Launch the server
	log.Fatal(server.ListenAndServe())
}



