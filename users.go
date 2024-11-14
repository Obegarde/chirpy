package main

import(
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
	"github.com/obegarde/chirpy/internal/auth"
	"github.com/obegarde/chirpy/internal/database"
	"database/sql"
	
)


type JSONUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token	string	`json:"token"`
	RefreshToken string `json:"refresh_token"`
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
		token,err := auth.MakeJWT(userByEmail.ID,cfg.secret)
		if err != nil{
			respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access", err)
			return
	}
		refreshString, err := auth.MakeRefreshToken()
		if err != nil{
			respondWithError(w,http.StatusInternalServerError,"500 Could not generate hexstring", err)
			return
	}
		timeNow := time.Now().UTC()
		refreshToken, err := cfg.db.CreateRefreshToken(r.Context(),database.CreateRefreshTokenParams{
						Token:	refreshString,
						UserID: userByEmail.ID,
						ExpiresAt: timeNow.Add(time.Hour*24*60),
						RevokedAt: sql.NullTime{
			Time:time.Time{},
			Valid:false,
		},
						
	})
	if err != nil{
		respondWithError(w,http.StatusInternalServerError,"Failed to make refresh token", err)
		return
	}

	respondWithJSON(w,http.StatusOK,JSONUser{
			ID:userByEmail.ID,
			CreatedAt:userByEmail.CreatedAt,
			UpdatedAt:userByEmail.UpdatedAt,
			Email:userByEmail.Email,
			Token:token,
			RefreshToken:refreshToken.Token,
	})

}

func (cfg *apiConfig)handlerCheckRefresh(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Auth header: %q\n", r.Header.Get("Authorization"))
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
    	fmt.Printf("Error: %v\n", err)
    		respondWithError(w,http.StatusBadRequest,"Failed to read header",err)
		return
	}
	fmt.Printf("Token: %q\n", tokenString)
	refreshToken, err:= cfg.db.GetRefreshToken(r.Context(),tokenString)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	if refreshToken.RevokedAt.Valid{		
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",fmt.Errorf("Token has been revoked"))
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()){	
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",fmt.Errorf("Token has expired"))
		return
	}
	userID,err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	accessToken,err := auth.MakeJWT(userID, cfg.secret)	
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	type tokenResponse struct{
		Token string `json:"token"`
	}
	
	respondWithJSON(w, http.StatusOK,tokenResponse{
		Token:accessToken,
		})


}

func (cfg *apiConfig)handlerRevokeRefresh(w http.ResponseWriter, r *http.Request){
tokenString, err := auth.GetBearerToken(r.Header)	
	if err != nil{
		respondWithError(w,http.StatusBadRequest,"Failed to read header",err)
		return
	}
	refreshToken, err:= cfg.db.GetRefreshToken(r.Context(),tokenString)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	if refreshToken.RevokedAt.Valid{		
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()){	
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	respondWithError(w,http.StatusNoContent,"204 Token Revoked",err)
}
