package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"info441-finalproj/servers/gateway/models/users"
	"info441-finalproj/servers/gateway/sessions"
)

func (handler *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Request body must be JSON", http.StatusUnsupportedMediaType)
		return
	}
	newUser := &users.NewUser{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, newUser)
	user, err2 := newUser.ToUser()
	if err2 != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	userRes, insertErr := handler.UserStore.Insert(user)
	if insertErr != nil {
		fmt.Println(insertErr)
		http.Error(w, "User could not be added to the database", http.StatusBadRequest)
		return
	}
	state := &SessionState{time.Now(), *userRes}
	_, sessionErr := sessions.BeginSession(handler.SessionKey, handler.SessionStore, state, w)
	if sessionErr != nil {
		http.Error(w, "Session could not be established", http.StatusInternalServerError)
		return
	}
	userJSON, _ := json.Marshal(userRes)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}

func (handler *Handler) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, sessionErr := sessions.GetState(r, handler.SessionKey, handler.SessionStore, state)
	if sessionErr != nil {
		http.Error(w, "Current user is not authenticated", http.StatusUnauthorized)
		return
	}
	idString := path.Base(r.URL.Path)
	var id int64
	if r.Method == "GET" {
		if idString == "me" {
			id = state.User.ID
		} else {
			temp, idErr := strconv.Atoi(idString)
			if idErr != nil {
				http.Error(w, "Passed ID was not a valid ID", http.StatusBadRequest)
				return
			}
			id = int64(temp)
		}
		user, userErr := handler.UserStore.GetByID(id)
		if userErr != nil {
			http.Error(w, "User with passed ID was not found", http.StatusNotFound)
			return
		}
		userJSON, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)
	} else if r.Method == "PATCH" {
		if idString != "me" {
			temp, idErr := strconv.Atoi(idString)
			if idErr != nil {
				http.Error(w, "Passed ID was not a valid ID", http.StatusBadRequest)
				return
			}
			id = int64(temp)
			if id != state.User.ID {
				http.Error(w, "Cannot PATCH user data for non-authenticated user", http.StatusForbidden)
				return
			}
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Request body must be JSON", http.StatusUnsupportedMediaType)
			return
		}
		update := &users.Updates{}
		data, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(data, update)
		updateErr := state.User.ApplyUpdates(update)
		if updateErr != nil {
			http.Error(w, "User could not be updated", http.StatusBadRequest)
			return
		}
		userJSON, _ := json.Marshal(state.User)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)
	} else {
		http.Error(w, "Method must be PATCH or GET", http.StatusMethodNotAllowed)
		return
	}
}

func (handler *Handler) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Request body must be JSON", http.StatusUnsupportedMediaType)
		return
	}
	creds := &users.Credentials{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, creds)
	user, userErr := handler.UserStore.GetByEmail(creds.Email)
	if userErr != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	authErr := user.Authenticate(creds.Password)
	if authErr != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	state := &SessionState{}
	state.Time = time.Now()
	state.User = *user
	ip := r.RemoteAddr
	if len(r.Header.Get("X-Forwarded-For")) != 0 {
		ips := strings.Split(ip, ", ")
		ip = ips[0]
	}
	handler.UserStore.TrackLogin(user.ID, ip, state.Time)
	sessions.BeginSession(handler.SessionKey, handler.SessionStore, state, w)
	userJSON, _ := json.Marshal(state.User)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}

func (handler *Handler) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method must be DELETE", http.StatusMethodNotAllowed)
		return
	}
	if path.Base(r.URL.Path) != "mine" {
		http.Error(w, "Cannot end someone else's session!", http.StatusForbidden)
		return
	}
	_, sessionErr := sessions.EndSession(r, handler.SessionKey, handler.SessionStore)
	if sessionErr != nil {
		http.Error(w, "Could not end session", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("signed out"))
}
