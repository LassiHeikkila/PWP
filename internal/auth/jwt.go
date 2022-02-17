package auth

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

const (
	IssuerName = "taskey-auth-service"
)

type Controller interface{}

type authController struct {
	key []byte
}

type claims struct {
	User         string `json:"user"`
	Organization string `json:"organization"`

	jwt.StandardClaims
}

func NewController(key []byte) *authController {
	if key == nil {
		const desiredLength = 256
		key = make([]byte, desiredLength)
		n, err := rand.Read(key)
		if err != nil {
			log.Println("error generating key:", err)
			return nil
		}
		if n != desiredLength {
			log.Println("generated key too short:", n)
			return nil
		}
	}

	return &authController{
		key: key,
	}
}

func CreateClaims(user string, organization string) jwt.Claims {
	return &claims{
		User:         user,
		Organization: organization,
		StandardClaims: jwt.StandardClaims{
			Issuer: IssuerName,
		},
	}
}

func (a *authController) CreateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(a.key)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (a *authController) ValidateToken(tokenString string, user string, organization string) bool {
	// implementation inspired by example at https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.key, nil
	})
	if err != nil {
		log.Println("error parsing token:", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimedUser := claims["user"]
		claimerOrganization := claims["organization"]

		return claimedUser == user && claimerOrganization == organization
	}
	return false
}
