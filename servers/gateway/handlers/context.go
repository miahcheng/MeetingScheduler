package handlers

import (
	"info441-finalproj/servers/gateway/models/users"
	"info441-finalproj/servers/gateway/sessions"
)

// HandlerContext holds and shares the SigningKey,
// SessionStore and UserStore to share the values with
// the handler
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}

// NewHandlerContext blah
func NewHandlerContext(signingKey string, sessStore sessions.Store, useStore users.Store) *HandlerContext {
	return &HandlerContext{
		signingKey,
		sessStore,
		useStore,
	}
}
