package login

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"project/userlogin/connection"
	"project/userlogin/library"
	"project/userlogin/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conn := connection.Client
	Users := conn.Database("test").Collection("users")   //connection to users collection
	Errors := conn.Database("test").Collection("errors") //connection to errors collection

	var user model.User
	//fetching payload
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Invalid payload : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid payload : " + err.Error()})
		return
	}

	// check empty fields
	if user.Username == "" || user.Firstname == "" || user.Lastname == "" || user.Email == "" || user.Password == "" {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "All fields are required",
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "All fields are required"})
		return
	}
	//check valid email
	if !library.IsValidEmail(user.Email) {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Provide valid email address",
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Provide valid email address"})
		return
	}
	// Check if email already exists
	existingUser := model.User{}
	filter := bson.M{"email": primitive.Regex{Pattern: user.Email, Options: "i"}}
	err = Users.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Email already exists",
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusInternalServerError,
			ErrorDesccription: "Error while checking existing email : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Error while checking existing email : " + err.Error()})
		return
	}
	// check if username already exists
	filter = bson.M{"username": primitive.Regex{Pattern: user.Username, Options: "i"}}
	err = Users.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Username already exists",
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Username already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusInternalServerError,
			ErrorDesccription: "Error while checking existing username : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Error while checking existing username : " + err.Error()})
		return
	}

	// Generate unique ID for user
	user.ID = primitive.NewObjectID()
	user.CreatedOn = time.Now() //current time
	// Insert user into database
	_, err = Users.InsertOne(context.Background(), user)
	if err != nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Register",
			ErrorCode:         http.StatusInternalServerError,
			ErrorDesccription: "Failed to register user : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Failed to register user : " + err.Error()})
		return
	}
	//final response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ErrorResponse{Message: "A verification mail has been sent to your registered mail."})
	return
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conn := connection.Client
	Users := conn.Database("test").Collection("users")   // connection to users collection
	Errors := conn.Database("test").Collection("errors") // connection to errors collection
	user := struct {
		Username string
		Password string
	}{}
	//fetching payload
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Login",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Invalid payload : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid payload : " + err.Error()})
		return
	}

	// Validate empty fields
	if user.Username == "" || user.Password == "" {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Login",
			ErrorCode:         http.StatusBadRequest,
			ErrorDesccription: "Username and password are required",
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Username and password are required"})
		return
	}

	// Check if user exists
	existingUser := model.User{}
	err = Users.FindOne(context.Background(), bson.M{"username": user.Username, "password": user.Password}).Decode(&existingUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "Login",
				ErrorCode:         http.StatusUnauthorized,
				ErrorDesccription: "Invalid username or password : " + err.Error(),
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Invalid username or password : " + err.Error()})
		} else {
			Errors.InsertOne(context.Background(), model.ErrorObject{
				ID:                primitive.NewObjectID(),
				ApiName:           "Login",
				ErrorCode:         http.StatusInternalServerError,
				ErrorDesccription: "Error logging in user : " + err.Error(),
				CreatedOn:         time.Now(),
			})
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Error logging in user : " + err.Error()})
		}
		return
	}

	// Generate JWT token
	token, err := library.GenerateJWT(existingUser)
	if err != nil {
		Errors.InsertOne(context.Background(), model.ErrorObject{
			ID:                primitive.NewObjectID(),
			ApiName:           "Login",
			ErrorCode:         http.StatusInternalServerError,
			ErrorDesccription: "Failed to generate authentication token : " + err.Error(),
			CreatedOn:         time.Now(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "Failed to generate authentication token : " + err.Error()})
		return
	}
	//final response
	response := model.TokenResponse{
		Token: token,
		User:  existingUser,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

// UserInfo fetch user information from the auth token.
func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//getting user object from request header extracted from jwt token
	user := r.Context().Value("user").(model.User)
	json.NewEncoder(w).Encode(user)
	return
}
