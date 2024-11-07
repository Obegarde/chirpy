package main
import (
	"net/http"
	"fmt"
	
)

func (cfg *apiConfig) resetHitHandler(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request){
	if cfg.platform != "dev"{
		respondWithError(w, 403,"FORBIDDEN",fmt.Errorf("403 FORBIDDEN"))
		return
	}else{


	err := cfg.db.ResetUsers(r.Context())
	if err != nil{	
		respondWithError(w, http.StatusInternalServerError,"Could not reset users",err)
			return
	}
	}

}
