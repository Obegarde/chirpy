
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
	w.Header().Set("Content-Type","text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(content))

}



func(cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)

	})
}