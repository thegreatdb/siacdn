package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

type createAccountForm struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	StripeToken string `json:"stripe_token"`
}

func (s *HTTPDServer) handleCreateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}
	if err = r.Body.Close(); err != nil {
		s.JsonErr(w, "Could not read data: "+err.Error())
		return
	}

	var form createAccountForm
	if err = json.Unmarshal(data, &form); err != nil {
		s.JsonErr(w, "Could not decode JSON: "+err.Error())
		return
	}
	if form.Username == "" || len(form.Username) < 4 {
		s.JsonErr(w, "Invalid username (must be at least 3 characters long)")
		return
	}
	if form.Password == "" || len(form.Password) < 6 {
		s.JsonErr(w, "Invalid password (must be at least 5 characters long)")
		return
	}
	if form.StripeToken == "" {
		s.JsonErr(w, "Must include Stripe token to register")
		return
	}

	// TODO: Stripe validation
	// TODO: Add Stripe card or customer id to model and add it here

	acc, err := models.NewAccount(form.Username, form.Password)
	if err != nil {
		s.JsonErr(w, "Could not create new account"+err.Error())
		return
	}
	if err = s.db.SaveAccount(acc); err != nil {
		s.JsonErr(w, "Could not save created account"+err.Error())
		return
	}

	authToken, err := models.NewAuthToken(acc.ID)
	if err != nil {
		s.JsonErr(w, "Could not create new auth token: "+err.Error())
		return
	}
	if err = s.db.SaveAuthToken(authToken); err != nil {
		s.JsonErr(w, "Could not save created auth token"+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"account": acc, "auth_token": authToken})
}
