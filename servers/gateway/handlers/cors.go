package handlers

import (
	"net/http"
)

const accessControlAllowOrigin = "Access-Control-Allow-Origin"
const accessControlAllowMethods = "Access-Control-Allow-Methods"
const accessControlAllowHeaders = "Access-Control-Allow-Headers"
const accessControlExposedHeaders = "Access-Control-Expose-Headers"
const accessControlMaxAge = "Access-Control-Max-Age"

const allowedMethods = "GET, PUT, POST, PATCH, DELETE"
const allowedHeaders = "Content-Type, Authorization"
const exposedHeaders = "Authorization"
const maxAge = "600"

// CorsMW blah
type CorsMW struct {
	Handler http.Handler
}

func (c *CorsMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(accessControlAllowOrigin, "*")
	w.Header().Set(accessControlAllowMethods, allowedMethods)
	w.Header().Set(accessControlAllowHeaders, allowedHeaders)
	w.Header().Set(accessControlExposedHeaders, exposedHeaders)
	w.Header().Set(accessControlMaxAge, maxAge)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	c.Handler.ServeHTTP(w, r)
}

// NewCorsMW blah
func NewCorsMW(handlerToWrap http.Handler) *CorsMW {
	return &CorsMW{handlerToWrap}
}
