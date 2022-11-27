package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.uber.org/zap/zapcore"
	"github.com/dgrijalva/jwt-go"
	
	"gitlab.com/m0ta/lts/app/config"
)

// User is a JSON-Postgre user
type User struct {
	bun.BaseModel 					`bun:"table:users,alias:u"`
	ID                	uuid.UUID	`json:"id" bun:"id,pk"`
	Email             	string 		`json:"email" bun:"email"`
	Password          	string 		`json:"password,omitempty" bun:"-"`
	EncryptedPassword 	string 		`json:"-" bun:"encrypted_password"`
	NickName 			string		`json:"nickname" bun:"nickname"`
	CreatedAt 			time.Time 	`json:"created_at,omitempty" bun:"created_at"`
	LoginedAt 			time.Time 	`json:"logined_at,omitempty" bun:"logined_at"`
	UpdatedAt 			time.Time 	`json:"updated_at,omitempty" bun:"updated_at"`
	Token    			string 		`json:"token,omitempty" bun:"-"`
}

//Users ...
type Users []User

func (user *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("ID", 		user.ID.String())
	enc.AddString("Email", 		user.Email)
	enc.AddString("Password", 	user.Password)
	enc.AddString("EncryptedPassword", 	user.EncryptedPassword)
	enc.AddString("NickName", 	user.NickName)
	enc.AddString("Token", 		user.Token)
	enc.AddTime("CreatedAt",	user.CreatedAt)
	enc.AddTime("UpdatedAt",	user.UpdatedAt)
    return nil
}

// Token generates token for web session
func (user *User)GenerateToken() (string, error) {
	cfg := config.Get()

	value, err := time.ParseDuration(cfg.SessionDuration)
	if err != nil {
		return "", err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  		user.ID,
		"expired": 	time.Now().Add(value).Unix(),
	})
	return jwt.SignedString([]byte(cfg.SignedToken))
}