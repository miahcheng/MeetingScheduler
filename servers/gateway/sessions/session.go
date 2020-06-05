package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, ErrNoSessionID
	}
	newID, _ := NewSessionID(signingKey)
	store.Save(newID, sessionState)
	tempName := schemeBearer + newID.String()
	w.Header().Add(headerAuthorization, tempName)
	return newID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	auth := r.Header.Get(headerAuthorization)
	if len(auth) == 0 {
		auth = r.URL.Query().Get("auth")
	}
	if !strings.Contains(auth, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	} else {
		auth = strings.TrimPrefix(auth, schemeBearer)
	}
	newID, err := ValidateID(auth, signingKey)
	return newID, err
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	newID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	storeErr := store.Get(newID, sessionState)
	if storeErr != nil {
		return InvalidSessionID, storeErr
	}
	return newID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	newID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	store.Delete(newID)
	return newID, nil
}
