package authApp

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nickname     string             `json:"nickname,omitempty" bson:"nickname,omitempty"`
	PasswordHash string             `json:"-" bson:"password,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
}

type CreateUser struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}
