package main

import(
	"net/http"
	"encoding/json"
	"fmt"
)



func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
	}
		
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err  != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest,"Could not decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email) 
	if err != nil{
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError,"Could not create user",err)
		return
	}
	respondWithJSON(w, 201, user) 
	
}

