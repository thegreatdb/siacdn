package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
	"github.com/thegreatdb/siacdn/backend/models"
)

func (s *HTTPDServer) handleDeleteSiaNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	admin := r.URL.Query().Get("secret") == SiaCDNSecretKey
	var account *models.Account
	var err error
	if !admin {
		req := NewRequest(r, s.db)
		account, err = req.GetAccount()
		if err != nil && err != db.ErrNotFound {
			s.JsonErr(w, err.Error())
			return
		}
		if account == nil || err == db.ErrNotFound {
			s.JsonErr(w, "You must be authenticated to access this resource")
			return
		}
	}

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

	if account != nil && sn.AccountID != account.ID {
		s.JsonErr(w, "You must be authenticated under the correct account to modify this resource")
		return
	}

	sn.Status = models.SIANODE_STATUS_STOPPING

	if err = s.db.SaveSiaNode(sn); err != nil {
		s.JsonErr(w, "Could not save Sia node: "+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"status": "ok"})
}
