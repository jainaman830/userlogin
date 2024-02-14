package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User : user object for registration
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Firstname string             `json:"firstname" bson:"firstname"`
	Lastname  string             `json:"lastname" bson:"lastname"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	CreatedOn time.Time          `json:"CreatedOn" bson:"CreatedOn"`
}

// TokenResponse : login output
type TokenResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ErrorResponse : error object
type ErrorResponse struct {
	Message string `json:"message"`
}

// ErrorObject : store error details
type ErrorObject struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	ApiName           string             `json:"ApiName" bson:"ApiName"`
	ErrorCode         int                `json:"ErrorCode" bson:"ErrorCode"`
	ErrorDesccription string             `json:"ErrorDesccription" bson:"ErrorDesccription"`
	UserID            string             `json:"UserID" bson:"UserID"`
	CreatedOn         time.Time          `json:"CreatedOn" bson:"CreatedOn"`
}
