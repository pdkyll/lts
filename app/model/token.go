package model

import (
	"time"

	"go.uber.org/zap/zapcore"
	"github.com/uptrace/bun"
	"github.com/dgrijalva/jwt-go"
	
	"gitlab.com/m0ta/lts/app/config"
)

// Token is a JSON-Postgre user
type Token struct {
	bun.BaseModel 			`bun:"table:tokens,alias:t"`
	ID         	uint64		`json:"id" bun:"id,pk"`
	Domain      string 		`json:"domain" bun:"domain"`
	Token       string 		`json:"token" bun:"token"`
	CreatedAt 	time.Time 	`json:"created_at" bun:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at" bun:"updated_at"`
	RequestedAt	time.Time 	`json:"requested_at" bun:"requested_at"`
	ExpiredAt 	time.Time 	`json:"expired_at" bun:"expired_at"`
}

// ValidData is a JSON-Postgre user
type ValidData struct {
	Valid	bool 	`json:"valid"`
	Token   *Token 	`json:"token"`
}

//Tokens ...
type Tokens []Token

func (token *Token) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddUint64("ID", 	token.ID)
	enc.AddString("Domain",	token.Domain)
	enc.AddString("Token", 	token.Token)
	enc.AddTime("CreatedAt",	token.CreatedAt)
	enc.AddTime("UpdatedAt",	token.UpdatedAt)
	enc.AddTime("RequestedAt",	token.RequestedAt)
	enc.AddTime("ExpiredAt",	token.ExpiredAt)
    return nil
}

// Token generates token for web session
func (token *Token)GenerateToken() error {
	cfg := config.Get()

	value, err := time.ParseDuration(cfg.TokenDuration)
	if err != nil {
		return err
	}
	
	token.ExpiredAt = token.CreatedAt.Add(value)

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  			token.ID,
		"domain":  		token.Domain,
		"expired_at": 	token.ExpiredAt,
	})

	token.Token, err = jwt.SignedString([]byte(cfg.SignedToken))
	if err != nil {
		return err
	}
	
	return nil
}