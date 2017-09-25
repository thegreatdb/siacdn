package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/siacdn-backend/db"
)

func (s *HTTPDServer) handleGetAuth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := NewRequest(r, s.db)
	account, err := req.GetAccount()
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"account": account})
}
