package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/my/repo/models/users"
	"github.com/my/repo/sessions"
)

// UsersHandler handles requestions for users to
// POST to create a new user account
func (context HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	reqMeth := r.Method
	if reqMeth != "POST" {
		http.Error(w, "Incorrect Status Method", http.StatusMethodNotAllowed)
		return
	}
	header := r.Header.Get("Content-Type")

	if !strings.HasPrefix(header, "application/json") {
		http.Error(w, "Header request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}
	dec := json.NewDecoder(r.Body)
	curUser := &users.NewUser{}
	dec.Decode(curUser)
	if valErr := curUser.Validate(); valErr != nil {
		// fmt.Printf("Error, Invalid user: %v", valErr)
		http.Error(w, "Error message: "+valErr.Error(), 400)
		return
	}
	toValUser, toUserErr := curUser.ToUser()
	if toUserErr != nil {
		fmt.Printf("Error creating a User: %v", toUserErr)
	}
	// insert new user into database
	cur, insertErr := context.UserStore.Insert(toValUser)
	if cur != toValUser || insertErr != nil {
		fmt.Printf("Error inserting new user into database: %v", insertErr)
		http.Error(w, "Error inserting user into database"+insertErr.Error(), 400)
	}
	// begin session
	sessState := &SessionState{BeginTime: time.Now(), CurrentUser: toValUser}
	// sessID, _ :=
	_, begErr := sessions.BeginSession(context.SigningKey, context.SessionStore, sessState, w)
	if begErr != nil {
		http.Error(w, "Error beginning session: "+begErr.Error(), http.StatusInternalServerError)
	}
	// context.SessionStore.Save(sessID, sessState)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if toValUser.ID != cur.ID {
		fmt.Printf("Error incorrect primary keys")
	}
	response, encErr := json.Marshal(toValUser)
	if encErr != nil {
		fmt.Printf("error encoding user to JSON: %v", encErr)
	}
	w.Write(response)
}

// SpecificUserHandler blah
func (context HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "PATCH" {
		http.Error(w, "Status Method Not Get or Patch", http.StatusMethodNotAllowed)
		return
	}
	idURL := path.Base(r.URL.Path)
	state := &SessionState{}
	var usrID int64
	_, err := sessions.GetState(r, context.SigningKey, context.SessionStore, state)
	if err != nil {
		http.Error(w, "error getting state: "+err.Error(), 401)
		return
	}
	if idURL == "me" {
		usrID = state.CurrentUser.ID
	} else {
		usrID, err = strconv.ParseInt(idURL, 10, 64)
		if err != nil {
			http.Error(w, "error parsing id "+err.Error(), 403)
			return
		}
	}
	curUser, getErr := context.UserStore.GetByID(usrID)
	if getErr != nil {
		http.Error(w, "Error getting user "+err.Error(), http.StatusInternalServerError)
	}
	// Get method
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(curUser)
		w.Write(response)
	} else if r.Method == "PATCH" {
		if idURL != "me" && usrID != state.CurrentUser.ID {
			http.Error(w, "Users do not match", http.StatusForbidden)
		}
		reqBody := r.Header.Get("Content-Type")
		if !strings.HasPrefix(reqBody, "application/json") {
			http.Error(w, "Request Body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		dec := json.NewDecoder(r.Body)
		updates := &users.Updates{}
		decErr := dec.Decode(updates)
		if decErr != nil {
			fmt.Printf("Error decoding json: %v", decErr)
		}
		// curUser := state.CurrentUser
		upUser, upEr := context.UserStore.Update(usrID, updates)
		if upEr != nil {
			http.Error(w, "Error updating user: "+upEr.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp, encErr := json.Marshal(upUser)
		if encErr != nil {
			fmt.Printf("Error encoding: %v", encErr)
		}
		w.Write(resp)

	} else {
		http.Error(w, "Status Method Not Get or Patch", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler blah
func (context HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reqBody := r.Header.Get("Content-Type") // im dumb
		if !strings.HasPrefix(reqBody, "application/json") {
			http.Error(w, "Requestion Body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		creds := &users.Credentials{}
		dec := json.NewDecoder(r.Body)
		decErr := dec.Decode(creds)
		if decErr != nil {
			fmt.Printf("Error decoding json: %v", decErr)
		}
		curUser, getErr := context.UserStore.GetByEmail(creds.Email)
		if getErr != nil {
			time.Sleep(time.Second)
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("User %s attempted to sign in", curUser.UserName)
		signInTime := time.Now()
		ipAdd := ""
		if len(r.Header.Get("X-Forwarded-For")) != 0 {
			ipAdd = r.Header.Get("X-Forwarded-For")
		} else {
			ipAdd = r.RemoteAddr
		}
		context.UserStore.InsertSignIn(curUser, signInTime, ipAdd)
		authErr := curUser.Authenticate(creds.Password)
		if authErr != nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		state := &SessionState{signInTime, curUser}
		sessions.BeginSession(context.SigningKey, context.SessionStore, state, w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp, encErr := json.Marshal(curUser)
		if encErr != nil {
			fmt.Printf("Error encoding: %v", encErr)
		}
		w.Write(resp)
	} else {
		http.Error(w, "Status Method is not Post", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificSessionHandler balh
func (context HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		path := r.URL.Path
		if !strings.HasSuffix(path, "mine") {
			http.Error(w, "Incorrect path to user", http.StatusForbidden)
			return
		}
		_, endErr := sessions.EndSession(r, context.SigningKey, context.SessionStore)
		if endErr != nil {
			http.Error(w, "Error ending sessions: "+endErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "Incorrect Status Method", http.StatusMethodNotAllowed)
		return
	}
}
