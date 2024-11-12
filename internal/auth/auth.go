package auth
import(	
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
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


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	claimsStruct := jwt.RegisteredClaims{
		Issuer:"chirpy",
		IssuedAt:jwt.NewNumericDate(time.Now()),
		ExpiresAt:jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256,claimsStruct)
	return newToken.SignedString(tokenSecret)

}
