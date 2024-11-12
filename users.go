package main

import(
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
	"github.com/obegarde/chirpy/internal/auth"
	"github.com/obegarde/chirpy/internal/database"
	
)


type JSONUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
		
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err  != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest,"Could not decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Falied to hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			Email:params.Email,
			HashedPassword: hashed_password,
	}) 
	if err != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError,"Could not create user",err)
		return
	}
	respondWithJSON(w, 201, JSONUser{
				ID: user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:	user.Email,
	}) 
	
}

func(cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest,"Could not decode parameters", err)
		return
	}
	userByEmail, err := cfg.db.GetUserByEmail(r.Context(),params.Email)	
	if err != nil{
		respondWithError(w, http.StatusUnauthorized,"Incorrect email or password", err)
		fmt.Println(err)
		return
	}	
	err = auth.CheckPasswordHash(params.Password, userByEmail.HashedPassword)
	if err != nil{ 
		respondWithError(w, http.StatusUnauthorized,"Incorrect email or password", fmt.Errorf("Incorrect email or password"))
		return
	}
	respondWithJSON(w,http.StatusOK,JSONUser{
			ID:userByEmail.ID,
			CreatedAt:userByEmail.CreatedAt,
			UpdatedAt:userByEmail.UpdatedAt,
			Email:userByEmail.Email,
	})

}
