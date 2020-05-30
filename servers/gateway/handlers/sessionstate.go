package handlers

import (
	"time"

	"github.com/my/repo/models/users"
)

// SessionState instantiates the time the Current User begins
// their session
type SessionState struct {
	BeginTime   time.Time
	CurrentUser *users.User
}
