package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (s *HTTPDServer) handleGetSiaNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		s.JsonErr(w, "Invalid (UUID) SiaNode ID: "+ps.ByName("id"))
		return
	}

	sn, err := s.db.GetSiaNode(id)
	if err != nil {
		s.JsonErr(w, "Could not get Sia node: "+err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"sianode": sn})
}
