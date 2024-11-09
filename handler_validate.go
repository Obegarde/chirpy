package main

import(	
	"strings"
	"slices"	
	"fmt"
)



func ValidateChirp(chirpParams JSONChirpParams)(bool, JSONChirpParams){
	const maxChirpLength = 140
	if len(chirpParams.Body) > maxChirpLength{	
		return false,chirpParams
	}
	swearWords := []string{"kerfuffle","sharbert","fornax"}

	splitWords := strings.Split(chirpParams.Body," ")

	for index, word := range splitWords{
		if slices.Contains(swearWords, strings.ToLower(word)){
			splitWords[index] = "****"
		}
	}
	validatedString := strings.Join(splitWords," ")
	chirpParams.Body = validatedString
	fmt.Printf("chirpPrams at end of chirp validation: %v\n",chirpParams)
		return true, chirpParams 
	}



