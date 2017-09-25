package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/db"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

type Request struct {
	db                *db.Database
	authTokenIDHeader string
	authToken         *models.AuthToken
	account           *models.Account
	hasCheckedAuth    bool
}

func NewRequest(r *http.Request, database *db.Database) *Request {
	req := &Request{
		db:                database,
		authTokenIDHeader: r.Header.Get("X-Auth-Token-ID"),
	}
	return req
}

func (r *Request) checkAuth() error {
	if r.authTokenIDHeader != "" {
		authTokenID, err := uuid.Parse(r.authTokenIDHeader)
		if err != nil {
			return err
		}
		authToken, err := r.db.GetAuthToken(authTokenID)
		if err != nil {
			return err
		}
		if authToken == nil {
			return nil
		}
		account, err := r.db.GetAccount(authToken.AccountID)
		if err != nil {
			return err
		}
		if account != nil {
			r.account = account
		}
	}
	return nil
}

func (r *Request) GetAccount() (*models.Account, error) {
	if !r.hasCheckedAuth {
		r.hasCheckedAuth = true
		if err := r.checkAuth(); err != nil {
			return nil, err
		}
	}
	if r.account == nil {
		return nil, errors.New("Account not found")
	}
	return r.account, nil
}
