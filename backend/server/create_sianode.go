package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/backend/db"
	"github.com/thegreatdb/siacdn/backend/models"
)

type createSiaNodeForm struct {
	Capacity float32 `json:"capacity"`
}

func (s *HTTPDServer) handleCreateSiaNode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}
	if err = r.Body.Close(); err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}

	var form createSiaNodeForm
	if err = json.Unmarshal(data, &form); err != nil {
		s.JsonErr(w, "Could not decode JSON: "+err.Error())
		return
	}
	if form.Capacity < 5 || form.Capacity > 50 {
		s.JsonErr(w, "Invalid capacity")
		return
	}

	sn, err := models.NewSiaNode(account.ID, form.Capacity)
	if err != nil {
		s.JsonErr(w, "Could not create new Sia node: "+err.Error())
		return
	}

	// TODO: Look up and make sure a Sia node with the same shortcode doesn't
	//       already exist, and if so, generate a new shortcode.

	if err = s.db.SaveSiaNode(sn); err != nil {
		s.JsonErr(w, "Could not save created Sia node: "+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"sianode": sn})
}
