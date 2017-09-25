package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST = 10

type Account struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedTime  time.Time `json:"created_time"`
}

func NewAccount(username, password string) (*Account, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	acc := &Account{
		ID:          id,
		Username:    username,
		CreatedTime: time.Now().UTC(),
	}
	if err = acc.SetPassword(password); err != nil {
		return nil, err
	}
	return acc, nil
}

func (acc *Account) SetPassword(password string) error {
	hsh, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return err
	}
	acc.PasswordHash = string(hsh)
	return nil
}

func (acc *Account) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(acc.PasswordHash),
		[]byte(password),
	)
}
