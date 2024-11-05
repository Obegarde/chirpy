package main

import(
	"encoding/json"
	"net/http"
	"strings"
	"slices"
)



func handlerChirpsValidate(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}
	type returnVals struct{
		Valid bool `json:"valid"`
	}
	type validatedParameters struct{
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError,"Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	

	swearWords := []string{"kerfuffle","sharbert","fornax"}

	splitWords := strings.Split(params.Body," ")

	for index, word := range splitWords{
		if slices.Contains(swearWords, strings.ToLower(word)){
			splitWords[index] = "****"
		}
	}
	validatedString := strings.Join(splitWords," ")
		

	respondWithJSON(w,http.StatusOK, validatedParameters{
		Body:validatedString,
	})


}
