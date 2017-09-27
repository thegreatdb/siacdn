package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type updateSiaNodeStatusForm struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (s *HTTPDServer) handleUpdateSiaNodeStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("secret") != SiaCDNSecretKey {
		s.JsonErr(w, "Secret key must match")
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}
	if err = r.Body.Close(); err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}

	var form updateSiaNodeStatusForm
	if err = json.Unmarshal(data, &form); err != nil {
		s.JsonErr(w, "Could not decode JSON: "+err.Error())
		return
	}

	sn, err := s.db.GetSiaNode(form.ID)
	if err != nil {
		s.JsonErr(w, "Could not get Sia node: "+err.Error())
		return
	}

	sn.Status = form.Status
	if err = sn.ValidateStatus(); err != nil {
		s.JsonErr(w, err.Error())
		return
	}

	if err = s.db.SaveSiaNode(sn); err != nil {
		s.JsonErr(w, "Could not save Sia node: "+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"sianode": sn})
}
