package handlers

import (
	"info441-finalproj/servers/gateway/models/users"
	"time"
)

type SessionState struct {
	Time time.Time  `json:"time"`
	User users.User `json:"user"`
}
