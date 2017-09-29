package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
)

func (s *HTTPDServer) handleGetWalletSeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("secret") != SiaCDNSecretKey {
		s.JsonErr(w, "Secret key must match")
		return
	}
	siaNodeIDStr := ps.ByName("id")
	siaNodeID, err := uuid.Parse(siaNodeIDStr)
	if err != nil {
		s.JsonErr(w, "Could not parse id '"+siaNodeIDStr+"': "+err.Error())
		return
	}
	seed, err := s.db.GetWalletSeed(siaNodeID)
	if err != nil && err != db.ErrNotFound {
		s.JsonErr(w, "Could not get wallet seed: "+err.Error())
		return
	}
	s.Json(w, map[string]interface{}{"wallet_seed": seed})
}
