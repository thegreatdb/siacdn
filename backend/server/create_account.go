package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/thegreatdb/siacdn/siacdn-backend/db"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

var stripeSecretKey string

func getStripeSecretKey() string {
	if stripeSecretKey == "" {
		stripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
	}
	if stripeSecretKey == "" {
		stripeSecretKey = "sk_test_cfNIKOMGfV4gI9KxzPxv4oba" // Testnet secret key
	}
	return stripeSecretKey
}

type createAccountForm struct {
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Name        string        `json:"name"`
	StripeToken *stripe.Token `json:"stripe_token"`
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
	if form.Email == "" || len(form.Email) < 6 {
		s.JsonErr(w, "Invalid email (must be at least 5 characters long)")
		return
	}
	if form.Password == "" || len(form.Password) < 6 {
		s.JsonErr(w, "Invalid password (must be at least 5 characters long)")
		return
	}
	if form.StripeToken == nil {
		s.JsonErr(w, "Must include Stripe token to register")
		return
	}

	// Check if the account exists already
	_, err = s.db.GetAccountByEmail(form.Email)
	if err == nil {
		s.JsonErr(w, "Account with that e-mail address already exists")
		return
	} else {
		if err == db.ErrNotFound {
			// This was the expected path
		} else {
			s.JsonErr(w, "Could not validate account uniqueness: "+err.Error())
			return
		}
	}

	stripe.Key = getStripeSecretKey()
	customerParams := &stripe.CustomerParams{
		Desc:  "Customer for " + form.Email,
		Email: form.Email,
	}
	customerParams.SetSource(form.StripeToken.ID)
	cust, err := customer.New(customerParams)
	if err != nil {
		s.JsonErr(w, "Could not create Stripe customer: "+err.Error())
		return
	}

	acc, err := models.NewAccount(form.Email, form.Password, form.Name, cust.ID)
	if err != nil {
		s.JsonErr(w, "Could not create new account: "+err.Error())
		return
	}
	if err = s.db.SaveAccount(acc); err != nil {
		s.JsonErr(w, "Could not save created account: "+err.Error())
		return
	}

	authToken, err := models.NewAuthToken(acc.ID)
	if err != nil {
		s.JsonErr(w, "Could not create new auth token: "+err.Error())
		return
	}
	if err = s.db.SaveAuthToken(authToken); err != nil {
		s.JsonErr(w, "Could not save created auth token: "+err.Error())
		return
	}

	s.Json(w, map[string]interface{}{"account": acc, "auth_token": authToken})
}
