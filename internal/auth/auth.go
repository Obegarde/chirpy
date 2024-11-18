package auth
import(	
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
	"fmt"
	"net/http"
	"strings"
	"crypto/rand"
	"encoding/hex"
	
)


func HashPassword(password string)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil{
		return "",err
	}
	return string(hashedPassword), nil
}

//Returns nil on success and an error on fail
func CheckPasswordHash(password string, hash string) error{
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error){
	expiresIn := time.Hour
	now := time.Now().UTC()
	secretBytes := []byte(tokenSecret)
	claimsStruct := jwt.RegisteredClaims{
		Issuer:"chirpy",
		IssuedAt:jwt.NewNumericDate(now),
		ExpiresAt:jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256,claimsStruct)
	signedToken, err := newToken.SignedString(secretBytes)
	if err != nil{
		return "" , err
	}
	return signedToken,nil

}

func ValidateJWT(tokenString, tokenSecret string)(uuid.UUID, error){
	fmt.Printf("Received token bytes: %v\n", []byte(tokenString))	
	secretBytes := []byte(tokenSecret)
	CustomClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString,&CustomClaims, func(token *jwt.Token)(interface{},error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretBytes, nil
	})
	if err != nil{	
		return uuid.UUID{},err
			
	}
	if token.Valid{
		tokenUUID,err := uuid.Parse(CustomClaims.Subject)
		if err != nil {
			return uuid.UUID{},err	
		} 
		return tokenUUID, nil
	}
	return uuid.UUID{}, fmt.Errorf("Invalid token")
}

func GetBearerToken(headers http.Header)(string,error){
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == ""{
		return "", fmt.Errorf("No Authorization Header found")
	}
	if !strings.HasPrefix(authorizationHeader,"Bearer"){
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authorizationHeader,"Bearer")
	return strings.TrimSpace(token), nil
}

func MakeRefreshToken() (string,error){
	amountOfBytes := 32
	byteSlice := make([]byte,amountOfBytes)
	_,err := rand.Read(byteSlice)
	if err != nil{
		return "", err
	}
	hexString := hex.EncodeToString(byteSlice) 
	if hexString == ""{
		return "", fmt.Errorf("hexString creation failed") 
	}
		return hexString, nil	
}


func GetApiKey(headers http.Header)(string,error){
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == ""{
		return "", fmt.Errorf("No Authorization Header found")
	}
	if !strings.HasPrefix(authorizationHeader,"ApiKey"){
		return "", fmt.Errorf("invalid authorization header format")
	}

	key := strings.TrimPrefix(authorizationHeader,"ApiKey")
	return strings.TrimSpace(key), nil
}
