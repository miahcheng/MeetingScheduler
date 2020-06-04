package handlers

import (
	"time"

	"info441-finalproj/servers/gateway/models/users"
)

// SessionState instantiates the time the Current User begins
// their session
type SessionState struct {
	BeginTime   time.Time
	CurrentUser *users.User
}
