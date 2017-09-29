package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
)

func (s *HTTPDServer) handleGetPendingSiaNodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("secret") != SiaCDNSecretKey {
		s.JsonErr(w, "Secret key must match")
		return
	}
	sns, err := s.db.GetPendingSiaNodes()
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, "Could not get pending sia nodes: "+err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"sianodes": sns})
}
