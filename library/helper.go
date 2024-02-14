package library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"project/userlogin/connection"
	"project/userlogin/model"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWT secret key
var jwtSecret = []byte("qOXuhecyg2N01F1itNjIUPB7rqFeWvMd")

// GenerateJWT generates a JWT token for authentication.
func GenerateJWT(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID.Hex()
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Second * 60).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthenticateMiddleware is a middleware function to authenticate users using JWT token.
func AuthenticateMiddleware(next http.Handler) http.Handler {
	conn := connection.Client
	Errors := conn.Database("test").Collection("errors")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "AuthenticateMiddleware",
				ErrorCode:         http.StatusUnauthorized,
				ErrorDesccription: "Authorization header is missing",
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Authorization header is missing"})
			return
		}
		auth := strings.Split(authHeader, " ")
		if len(auth) != 2 {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "AuthenticateMiddleware",
				ErrorCode:         http.StatusUnauthorized,
				ErrorDesccription: "Invalid Authorization header",
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid Authorization header"})
			return
		}
		tokenString := auth[1] // Extract token from "Bearer <token>" format
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "AuthenticateMiddleware",
				ErrorCode:         http.StatusUnauthorized,
				ErrorDesccription: "Invalid token",
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user information from claims
			userid, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))
			user := model.User{
				ID:       userid,
				Username: fmt.Sprintf("%v", claims["username"]),
				Email:    fmt.Sprintf("%v", claims["email"]),
			}

			// Add user information to request context
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "AuthenticateMiddleware",
				ErrorCode:         http.StatusUnauthorized,
				ErrorDesccription: "Invalid token",
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid token"})
			return
		}
	})
}

func IsValidEmail(email string) bool {
	// Regular expression for validating email addresses
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
