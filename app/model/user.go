package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

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
	UpdatedAt 			time.Time 	`json:"updated_at,omitempty" bun:"updated_at"`
}

//Users ...
type Users []User