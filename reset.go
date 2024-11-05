package main
import "net/http"

func (cfg *apiConfig) resetHitHandler(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
}
