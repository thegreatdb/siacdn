package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
)

func (s *HTTPDServer) handleGetSiaNodesByAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	sns, err := s.db.GetSiaNodesByAccount(account.ID)
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, "Could not get sia nodes: "+err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"sianodes": sns})
}
