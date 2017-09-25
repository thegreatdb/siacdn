package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	CreatedTime time.Time `json:"created_time"`
}

func NewAuthToken(accountID uuid.UUID) (*AuthToken, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &AuthToken{
		ID:          id,
		AccountID:   accountID,
		CreatedTime: time.Now().UTC(),
	}, nil
}
