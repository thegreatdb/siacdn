package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedTime time.Time `json:"created_time"`
}

func NewAuthToken(userID uuid.UUID) (*AuthToken, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &AuthToken{
		ID:          id,
		UserID:      userID,
		CreatedTime: time.Now().UTC(),
	}, nil
}
