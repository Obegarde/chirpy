package main
import (
	"net/http"
	"encoding/json"
	"fmt"	
	"github.com/obegarde/chirpy/internal/database"
	"time"
	"github.com/google/uuid"
	"database/sql"
	"github.com/obegarde/chirpy/internal/auth"
)

type JSONChirp struct{
		ID        uuid.UUID 	`json:"id"`
		CreatedAt time.Time	`json:"created_at"`
		UpdatedAt time.Time	`json:"updated_at"`
		Body      string	`json:"body"`
		UserID    uuid.UUID	`json:"user_id"`
}
type JSONChirpParams struct {
    Body   string    `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}


func(cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){		
	decoder := json.NewDecoder(r.Body)
	params := JSONChirpParams{}
	err := decoder.Decode(&params)
	if err  != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest,"Could not decode parameters", err)
		return
	}
	tokenString,err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized,"401 Unauthorized",err)
		return
	}
	idFromToken, err := auth.ValidateJWT(tokenString,cfg.secret)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized,"401 Unauthorized", err)
		return
	}


	validatedUser, err := cfg.db.GetUserByID(r.Context(),idFromToken)
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"Could not validate user", fmt.Errorf("Could not validate user"))
		return
	}
	params.UserID = validatedUser.ID		
	validated, validatedParams := ValidateChirp(params)
	if !validated{
		respondWithError(w, http.StatusBadRequest,"Chirp too long", fmt.Errorf("Bad chirp"))
		return
	}else{
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body: validatedParams.Body,
			UserID: validatedParams.UserID,
			})
	if err != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError,"Could not create chirp",err)
		return
	}
	respondWithJSON(w, 201, JSONChirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
				
		})
	}
	
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request){
	allJSONChirps := []JSONChirp{}
	allChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Could not get chirps", err)
		return
	}
	for _, chirp := range allChirps{
		allJSONChirps = append(allJSONChirps,JSONChirp{
			ID:chirp.ID,
			CreatedAt:chirp.CreatedAt,
			UpdatedAt:chirp.UpdatedAt,
			Body:chirp.Body,
			UserID:chirp.UserID,
		})	
	}
	respondWithJSON(w,200,allJSONChirps)
				
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request){
	chirpID,err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest,"could not parse id",err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(),chirpID)
	if err == sql.ErrNoRows {
		respondWithError(w, 404,"chirp not found", fmt.Errorf("404 No chirp Found"))
		return
	}

	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Problem finding chirp", err)
		return
	}
	respondWithJSON(w,200,JSONChirp{
		ID:chirp.ID,
		CreatedAt:chirp.CreatedAt,
		UpdatedAt:chirp.UpdatedAt,
		Body:chirp.Body,
		UserID:chirp.UserID,
	})
	
	}

