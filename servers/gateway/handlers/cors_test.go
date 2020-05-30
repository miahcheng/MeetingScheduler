package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my/repo/models/users"
	"github.com/my/repo/sessions"
)

var newUserJ = []byte(`{"Email": "test@test.com", "Password": "password123", "PasswordConf": "password123", "UserName": "username", "FirstName": "firstname", "LastName": "lastname"}`)
var jsonStr = []byte(`{"ID": 1, "Email": "test@test.com", "PassHash": []byte("passhash123"), "UserName": "username", "FirstName": "firstname", "LastName": "lastname", "PhotoURL": "photourl"}`)
var upJSON = []byte(`{"FirstName": "hello", "LastName": "there"}`)

// TestMiddleUsersHandler tests the middleware handler and users handler
func TestMiddleUsersHandler(t *testing.T) {
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(newUserJ))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	custCtx := &HandlerContext{
		SigningKey:   "signingkey",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore: &users.FakeUserStore{
			User: &users.User{
				ID:        1,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			Err: nil,
		}}

	usersH := NewCorsMW(http.HandlerFunc(custCtx.UsersHandler))
	usersH.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("Expected %d got %d", http.StatusCreated, resp.Code)
	}
	if resp.Header().Get(accessControlAllowOrigin) != "*" {
		t.Errorf("Expected %s got %s", "*", resp.Header().Get(accessControlAllowOrigin))
	}
	if resp.Header().Get(accessControlAllowMethods) != allowedMethods {
		t.Errorf("Expected %s got %s", accessControlAllowMethods, resp.Header().Get(accessControlAllowMethods))
	}
	if resp.Header().Get(accessControlAllowHeaders) != allowedHeaders {
		t.Errorf("Expected %s got %s", allowedHeaders, resp.Header().Get(accessControlAllowHeaders))
	}
	if resp.Header().Get(accessControlExposedHeaders) != exposedHeaders {
		t.Errorf("Expected %s got %s", exposedHeaders, resp.Header().Get(accessControlExposedHeaders))
	}
	if resp.Header().Get(accessControlMaxAge) != maxAge {
		t.Errorf("Expected %s got %s", maxAge, resp.Header().Get(accessControlMaxAge))
	}

}

func TestMiddleSpecUserHandler(t *testing.T) {
	sessID, _ := sessions.NewSessionID("sessionid")
	newUse := &users.NewUser{
		Email:        "test@test.com",
		Password:     "passhash123",
		PasswordConf: "passhash123",
		UserName:     "username",
		FirstName:    "firstname",
		LastName:     "lastname",
	}
	usr, _ := newUse.ToUser()
	usr.ID = 1
	custCtx := &HandlerContext{
		SigningKey:   "sessionid",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore:    users.NewFakeUserStore(usr),
	}
	// fullSessID := "Bearer <" + sessID.String() + ">"
	fullSessID := "Bearer " + sessID.String()
	state := &SessionState{BeginTime: time.Now(), CurrentUser: usr}
	custCtx.SessionStore.Save(sessID, &state)
	req, _ := http.NewRequest("GET", "/v1/users/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fullSessID)
	resp := httptest.NewRecorder()
	middle := NewCorsMW(http.HandlerFunc(custCtx.SpecificUserHandler))
	middle.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusCreated, resp.Code)
	}
	if resp.Header().Get(accessControlAllowOrigin) != "*" {
		t.Errorf("Expected %s got %s", "*", resp.Header().Get(accessControlAllowOrigin))
	}
	if resp.Header().Get(accessControlAllowMethods) != allowedMethods {
		t.Errorf("Expected %s got %s", accessControlAllowMethods, resp.Header().Get(accessControlAllowMethods))
	}
	if resp.Header().Get(accessControlAllowHeaders) != allowedHeaders {
		t.Errorf("Expected %s got %s", allowedHeaders, resp.Header().Get(accessControlAllowHeaders))
	}
	if resp.Header().Get(accessControlExposedHeaders) != exposedHeaders {
		t.Errorf("Expected %s got %s", exposedHeaders, resp.Header().Get(accessControlExposedHeaders))
	}
	if resp.Header().Get(accessControlMaxAge) != maxAge {
		t.Errorf("Expected %s got %s", maxAge, resp.Header().Get(accessControlMaxAge))
	}
}
