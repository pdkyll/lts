package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	
	"gitlab.com/m0ta/lts/app/config"
	"gitlab.com/m0ta/lts/app/model"
)

//HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//VerifyPassword ...
func VerifyPassword(hash string, password string) error {
	bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// TokenGenerate generates the jwt token based on payload
func TokenGenerate(u *model.User) (string, error) {
	cfg := config.Get()
	

	v, err := time.ParseDuration(cfg.TokenExp)
	if err != nil {
		return "", err//panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(v).Unix(),
		"id":  u.ID,
	})

	token, err := t.SignedString([]byte(cfg.TokenKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func parse(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.Get().TokenKey), nil
	})
}

// TokenVerify verifies the jwt token against the secret
func TokenVerify(token string) (string, error) {
	parsed, err := parse(token)

	if err != nil {
		return "", err
	}

	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	// Getting ID, it's an interface{} so I need to cast it to uint
	id := claims["id"].(string)

	return id, nil
}