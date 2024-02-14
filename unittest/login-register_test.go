package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/userlogin/connection"
	"project/userlogin/login"
	"project/userlogin/model"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	// Initialize MongoDB connection for testing
	err := connection.ConnectDB()
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer connection.Client.Disconnect(context.Background())

	// Create a dummy user payload
	user := model.User{
		Username:  "jainaman",
		Firstname: "aman",
		Lastname:  "jain",
		Email:     "test@gmail.com",
		Password:  "password",
	}
	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Error marshalling user: %v", err)
	}

	// Create a POST request with the user payload
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request with the RegisterUser handler function
	handler := http.HandlerFunc(login.Register)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"A verification mail has been sent to your registered mail."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginUser(t *testing.T) {
	// Initialize MongoDB connection for testing
	err := connection.ConnectDB()
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer connection.Client.Disconnect(context.Background())

	// Create a dummy user payload
	user := model.User{
		Username: "jainaman",
		Password: "password",
	}
	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Error marshalling user: %v", err)
	}

	// Create a POST request with the user payload
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request with the LoginUser handler function
	handler := http.HandlerFunc(login.Login)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response model.TokenResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if response.Token == "" {
		t.Errorf("Handler returned empty token")
	}
}
