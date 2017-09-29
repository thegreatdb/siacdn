package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

type createWalletSeedForm struct {
	Words string `json:"words"`
}

func (s *HTTPDServer) handleCreateWalletSeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("secret") != SiaCDNSecretKey {
		s.JsonErr(w, "Secret key must match")
		return
	}
	siaNodeIDStr := ps.ByName("id")
	siaNodeID, err := uuid.Parse(siaNodeIDStr)
	if err != nil {
		s.JsonErr(w, "Invalid node ID: "+err.Error())
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

	var form createWalletSeedForm
	if err = json.Unmarshal(data, &form); err != nil {
		s.JsonErr(w, "Could not decode JSON: "+err.Error())
		return
	}

	seed, err := models.NewWalletSeed(siaNodeID, form.Words)
	if err != nil {
		s.JsonErr(w, "Could not create wallet seed: "+err.Error())
		return
	}

	if err = s.db.SaveWalletSeed(seed); err != nil {
		s.JsonErr(w, "Could not save wallet seed: "+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"wallet_seed": seed})
}
