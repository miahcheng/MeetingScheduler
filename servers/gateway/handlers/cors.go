package handlers

import (
	"net/http"
)

type Preflight struct {
	handler http.Handler
}

func (p *Preflight) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	p.handler.ServeHTTP(w, r)
}

func NewPreflight(handlerToWrap http.Handler) *Preflight {
	return &Preflight{handlerToWrap}
}
