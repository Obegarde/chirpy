package main
import(
	"sync/atomic"
	"database/sql"
	"github.com/obegarde/chirpy/internal/database"
	
)

type apiConfig struct{
	fileserverHits *atomic.Int32	
	db *database.Queries
}


func NewConfig(db *sql.DB) *apiConfig{
	var hits atomic.Int32
	return &apiConfig{
		fileserverHits: &hits,
		db : database.New(db),
	}
}
