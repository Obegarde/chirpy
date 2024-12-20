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
	IsChirpyRed bool `json:"is_chirpy_red"`
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
				IsChirpyRed: user.IsChirpyRed.Bool,
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
			IsChirpyRed: userByEmail.IsChirpyRed.Bool,
	})

}

func (cfg *apiConfig)handlerCheckRefresh(w http.ResponseWriter, r *http.Request){
	
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
    	fmt.Printf("Error: %v\n", err)
    		respondWithError(w,http.StatusBadRequest,"Failed to read header",err)
		return
	}
	
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


func (cfg *apiConfig)handlerUpdateUser(w http.ResponseWriter, r *http.Request){
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Failed to hash password",err)
		return
	}
	
	updateParams := database.UpdateUserParams{
		ID : idFromToken,
		Email : params.Email,
		HashedPassword : hashedPassword,
}
	updateResult,err := cfg.db.UpdateUser(r.Context(),updateParams)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Failed to update user information",err)
		return
	}
	type updateResponse struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`

	}
	respondWithJSON(w,http.StatusOK,updateResponse{
		ID:updateResult.ID,
		CreatedAt:updateResult.CreatedAt,
		UpdatedAt:updateResult.UpdatedAt,
		Email:updateResult.Email,
		IsChirpyRed:updateResult.IsChirpyRed.Bool,
	})
}

func (cfg *apiConfig)handlerUpgradeUser(w http.ResponseWriter,r *http.Request){
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized,"401 Unauthorized access",err)
		return
	}
	if apiKey != cfg.polkaKey{
		respondWithError(w,http.StatusUnauthorized,"401 Unauthorized access", fmt.Errorf("401 Unauthorized access"))
		return
	}
	type Data struct{
		UserID uuid.UUID `json:"user_id"`
	}
	type Event struct{
		Event string `json:"event"`
		Data Data `json:"data"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := Event{}
	err = decoder.Decode(&params)
	if err  != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest,"Could not decode parameters", err)
		return
	}
	if params.Event == "user.upgraded"{
		_,err := cfg.db.UpgradeUserToChirpyRed(r.Context(),params.Data.UserID)
		if err == sql.ErrNoRows{
			respondWithError(w, http.StatusNotFound,"Could not find User",err)	
			return
		}
		if err != nil{
			respondWithError(w,http.StatusInternalServerError,"Could not upgrade user",err)
			return
		}
		respondWithError(w,http.StatusNoContent,"",fmt.Errorf(""))
		return
	}else{
		respondWithError(w, http.StatusNoContent,"Not an upgrade event",fmt.Errorf("Only upgrade events to this hook"))	
		return
	}
	

	

}
