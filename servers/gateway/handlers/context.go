package handlers

import (
	"info441-finalproj/servers/gateway/models/users"
	"info441-finalproj/servers/gateway/sessions"
)

type Handler struct {
	SessionKey   string         `json:"key"`
	SessionStore sessions.Store `json:"sessions"`
	UserStore    users.Store    `json:"users"`
}
