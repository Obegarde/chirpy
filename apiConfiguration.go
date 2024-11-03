package main
import(
	"sync/atomic"
	"net/http"	
	"fmt"
)

type apiConfig struct{
	fileserverHits *atomic.Int32	
}


func NewConfig() *apiConfig{
	var hits atomic.Int32
	return &apiConfig{
		fileserverHits: &hits,
	}
}

func (cfg *apiConfig) hitHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(content))

}

func (cfg *apiConfig) resetHitHandler(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
}
