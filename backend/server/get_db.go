package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *HTTPDServer) handleGetDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("secret") != SiaCDNSecretKey {
		s.JsonErr(w, "Secret key must match")
		return
	}
	s.Json(w, s.db)
}
