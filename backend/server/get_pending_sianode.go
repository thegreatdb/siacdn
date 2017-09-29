package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
)

func (s *HTTPDServer) handleGetPendingSiaNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := NewRequest(r, s.db)
	account, err := req.GetAccount()
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, err.Error())
		return
	}
	if account == nil || err == db.ErrNotFound {
		s.JsonErr(w, "You must be authenticated to access this resource")
		return
	}
	sn, err := s.db.GetPendingSiaNode(account.ID)
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, "Could not get pending sia node: "+err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"sianode": sn})
}
