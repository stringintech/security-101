package test

import (
	"bytes"
	"encoding/json"
	"github.com/stringintech/security-101/model"
	"github.com/stringintech/security-101/server"
	"github.com/stringintech/security-101/server/auth"
	"github.com/stringintech/security-101/store"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTest() (*httptest.Server, *auth.JwtService) {
	userStore := store.NewUserStore()
	jwtService := auth.NewJwtService(auth.JwtServiceConfig{
		Secret:             []byte("test-secret"),
		ExpirationInterval: 10 * 24 * time.Hour,
	})

	return NewTestServer(userStore, jwtService), jwtService
}

func NewTestServer(userStore *store.UserStore, jwtService *auth.JwtService) *httptest.Server {
	s := server.New(userStore, jwtService)
	return httptest.NewServer(s)
}

func TestRegisterUser(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	reqBody := map[string]interface{}{
		"username":  "testuser",
		"password":  "Test123@pass",
		"full_name": "Test User",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if user.Username != "testuser" || user.FullName != "Test User" {
		t.Error("Response user details don't match request")
	}
}

func TestRegisterExistingUsername(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	reqBody := map[string]interface{}{
		"username":  "duplicate",
		"password":  "Test123@pass",
		"full_name": "Test User",
	}
	body, _ := json.Marshal(reqBody)

	// First registration
	_, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Duplicate registration
	resp, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest; got %v", resp.StatusCode)
	}
}

func TestLoginValidUser(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	// Register
	registerBody := map[string]interface{}{
		"username":  "logintest",
		"password":  "Test123@pass",
		"full_name": "Login Test",
	}
	body, _ := json.Marshal(registerBody)
	_, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Login
	loginBody := map[string]interface{}{
		"username": "logintest",
		"password": "Test123@pass",
	}
	body, _ = json.Marshal(loginBody)
	resp, err := http.Post(s.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}

	if loginResp["token"] == nil {
		t.Error("No token in response")
	}
}

func TestLoginWithNonexistentUsername(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	loginBody := map[string]interface{}{
		"username": "nonexistent",
		"password": "password",
	}
	body, _ := json.Marshal(loginBody)
	resp, err := http.Post(s.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized; got %v", resp.StatusCode)
	}
}

func TestLoginWithWrongPassword(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	// Register
	registerBody := map[string]interface{}{
		"username":  "wrongpass",
		"password":  "Test123@pass",
		"full_name": "Wrong Pass Test",
	}
	body, _ := json.Marshal(registerBody)
	_, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Login with wrong password
	loginBody := map[string]interface{}{
		"username": "wrongpass",
		"password": "wrongpassword",
	}
	body, _ = json.Marshal(loginBody)
	resp, err := http.Post(s.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized; got %v", resp.StatusCode)
	}
}

func TestGetUserWithValidToken(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	// Register and login to get token
	registerBody := map[string]interface{}{
		"username":  "getuser",
		"password":  "Test123@pass",
		"full_name": "Get User Test",
	}
	body, _ := json.Marshal(registerBody)
	_, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	loginBody := map[string]interface{}{
		"username": "getuser",
		"password": "Test123@pass",
	}
	body, _ = json.Marshal(loginBody)
	resp, err := http.Post(s.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	token := loginResp["token"].(string)

	// Get user details
	req, _ := http.NewRequest("POST", s.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if user.Username != "getuser" {
		t.Error("Response user details don't match")
	}
}

func TestGetUserWithoutAuthHeader(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	resp, err := http.Post(s.URL+"/users/me", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized; got %v", resp.StatusCode)
	}
}

func TestGetUserWithInvalidTokenFormat(t *testing.T) {
	s, _ := setupTest()
	defer s.Close()

	req, _ := http.NewRequest("POST", s.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.format")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized; got %v", resp.StatusCode)
	}
}

func TestGetUserWithExpiredToken(t *testing.T) {
	s, jwtService := setupTest()
	defer s.Close()

	// Override JWT service with shorter expiration
	*jwtService = *auth.NewJwtService(auth.JwtServiceConfig{
		Secret:             []byte("test-secret"),
		ExpirationInterval: 1 * time.Nanosecond,
	})

	// Register and login to get token
	registerBody := map[string]interface{}{
		"username":  "expired",
		"password":  "Test123@pass",
		"full_name": "Expired Token Test",
	}
	body, _ := json.Marshal(registerBody)
	_, err := http.Post(s.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	loginBody := map[string]interface{}{
		"username": "expired",
		"password": "Test123@pass",
	}
	body, _ = json.Marshal(loginBody)
	resp, err := http.Post(s.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	token := loginResp["token"].(string)

	// Wait for token to expire
	time.Sleep(2 * time.Nanosecond)

	// Try to get user details with expired token
	req, _ := http.NewRequest("POST", s.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized; got %v", resp.StatusCode)
	}
}
