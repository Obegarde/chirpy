package main
import(
	"sync/atomic"
	"database/sql"
	"github.com/obegarde/chirpy/internal/database"
	
)

type apiConfig struct{
	fileserverHits *atomic.Int32	
	db *database.Queries
	platform string
}


func NewConfig(db *sql.DB, currentplatform string) *apiConfig{
	var hits atomic.Int32
	return &apiConfig{
		fileserverHits: &hits,
		db : database.New(db),
		platform : currentplatform,
	}
}
