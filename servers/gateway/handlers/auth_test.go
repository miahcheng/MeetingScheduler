package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/my/repo/models/users"
	"github.com/my/repo/sessions"
)

// TestUsersHandler blah
func TestUsersHandler(t *testing.T) {
	var jsonStr = []byte(`{"Email": "test@test.com", "Password": "password123", "PasswordConf": "password123", "UserName": "username", "FirstName": "firstname", "LastName": "lastname"}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
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
	custCtx.UsersHandler(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("incorrect response status code: expected %d but got %d", http.StatusCreated, rr.Code)
	}
	expectedctype := "application/json"
	ctype := rr.Header().Get("Content-Type")
	if len(ctype) == 0 {
		t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
	} else if !strings.HasPrefix(ctype, expectedctype) {
		t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
	}

	postReq, _ := http.NewRequest("GET", "/", bytes.NewBuffer(jsonStr))
	postReq.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	custCtx.UsersHandler(resp, postReq)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Errorf("incorrect response status code, expected %d got %d", http.StatusMethodNotAllowed, resp.Code)
	}
	postReq2, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	postReq2.Header.Set("Content-Type", "plaintext")
	resp2 := httptest.NewRecorder()
	custCtx.UsersHandler(resp2, postReq2)
	if resp2.Code != http.StatusUnsupportedMediaType {
		t.Errorf("incorrect response status code, expected %d got %d", http.StatusUnsupportedMediaType, resp2.Code)
	}

}

func TestSpecificUserHandler(t *testing.T) {
	var jsonStr = []byte(`{"ID": 1, "Email": "test@test.com", "PassHash": []byte("passhash123"), "UserName": "username", "FirstName": "firstname", "LastName": "lastname", "PhotoURL": "photourl"}`)
	upJSON := []byte(`{"FirstName": "hello", "LastName": "there"}`)
	cases := []struct {
		name                string
		expectedUser        *users.User
		httpRequest         *http.Request
		expectedStatusCode  int
		expectedContentType string
		expectError         bool
	}{
		{
			"Valid Get Request Method",
			&users.User{
				ID:        1,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			httptest.NewRequest("GET", "/v1/users/1", bytes.NewBuffer(jsonStr)),
			200,
			"application/json",
			false,
		},
		{
			"Valid Patch Request Method",
			&users.User{
				ID:        1,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "hello",
				LastName:  "there",
				PhotoURL:  "photourl",
			},
			httptest.NewRequest("PATCH", "/v1/users/1", bytes.NewBuffer(upJSON)),
			200,
			"application/json",
			false,
		},
		{
			"Invalid Patch Request URL",
			&users.User{},
			httptest.NewRequest("PATCH", "/v1/users/you", bytes.NewBuffer(upJSON)),
			403,
			"",
			true,
		},
		{
			"Other HTTP Method not GET or PATCH",
			&users.User{},
			httptest.NewRequest("POST", "/v1/users/1", nil),
			405,
			"",
			true,
		},
		{
			"Invalid Patch Header",
			&users.User{},
			httptest.NewRequest("PATCH", "/v1/users/me", bytes.NewBuffer(upJSON)),
			415,
			"",
			true,
		},
	}
	sessID, _ := sessions.NewSessionID("sessionid")
	// fullSessID := "Bearer <" + sessID.String() + ">"
	fullSessID := "Bearer " + sessID.String()
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
	state := &SessionState{BeginTime: time.Now(), CurrentUser: usr}
	custCtx.SessionStore.Save(sessID, &state)

	for _, c := range cases {
		req := c.httpRequest
		if c.name == "Invalid Patch Header" {
			req.Header.Set("Content-Type", "plaintext")
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Authorization", fullSessID)
		resp := httptest.NewRecorder()
		custCtx.SpecificUserHandler(resp, req)
		log.Println(resp.Code)
		if resp.Code != c.expectedStatusCode {
			t.Errorf("Error %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatusCode, resp.Code)
		}
		if !c.expectError {
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", c.expectedContentType)
			} else if !strings.HasPrefix(ctype, c.expectedContentType) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", c.expectedContentType, ctype)
			}
		}

	}
}

func TestSessionsHandler(t *testing.T) {
	jsonStr := []byte(`{"Email": "test@test.com", "Password": "passhash123"}`)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	newUse := &users.NewUser{
		Email:        "test@test.com",
		Password:     "passhash123",
		PasswordConf: "passhash123",
		UserName:     "username",
		FirstName:    "firstname",
		LastName:     "lastname",
	}
	user, _ := newUse.ToUser()
	custCtx := &HandlerContext{
		SigningKey:   "signingkey",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore:    users.NewFakeUserStore(user),
	}
	custCtx.UserStore.Insert(user)
	custCtx.SessionsHandler(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("incorrect response status code: expected %d but got %d", http.StatusCreated, resp.Code)
	}
	expectedctype := "application/json"
	ctype := resp.Header().Get("Content-Type")
	if len(ctype) == 0 {
		t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
	} else if !strings.HasPrefix(ctype, expectedctype) {
		t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
	}

	req2, _ := http.NewRequest("GET", "/", nil)
	resp2 := httptest.NewRecorder()
	custCtx.SessionsHandler(resp2, req2)
	if resp2.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected %d got %d", http.StatusMethodNotAllowed, resp2.Code)
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	var jsonStr = []byte(`{"ID": 1, "Email": "test@test.com", "PassHash": []byte("passhash123"), "UserName": "username", "FirstName": "firstname", "LastName": "lastname", "PhotoURL": "photourl"}`)
	cases := []struct {
		name               string
		httpRequest        *http.Request
		expectedStatusCode int
		expectError        bool
	}{
		{
			"Valid Delete Request Method",
			httptest.NewRequest("DELETE", "/v1/sessions/mine", bytes.NewBuffer(jsonStr)),
			0,
			false,
		},
		{
			"Invalid, other HTTP Method",
			httptest.NewRequest("GET", "/v1/sessions/mine", nil),
			405,
			true,
		},
		{
			"Ivalid path segment",
			httptest.NewRequest("DELETE", "/v1/sessions/me", bytes.NewBuffer(jsonStr)),
			403,
			true,
		},
	}
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
	state := &SessionState{BeginTime: time.Now(), CurrentUser: usr}
	custCtx.SessionStore.Save(sessID, &state)
	for _, c := range cases {
		resp := httptest.NewRecorder()
		custCtx.SpecificSessionHandler(resp, c.httpRequest)
		if c.expectError {
			if resp.Code != c.expectedStatusCode {
				t.Errorf("Error %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatusCode, resp.Code)
			}
		}
	}
}
