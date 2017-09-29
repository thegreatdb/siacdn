package db

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/backend/models"
)

func (db *Database) GetAccount(id uuid.UUID) (*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if account, ok := db.Accounts[id]; ok {
		return account.Copy(), nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) GetAccountByEmail(email string) (*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, acc := range db.Accounts {
		if strings.ToLower(acc.Email) == strings.ToLower(email) {
			return acc.Copy(), nil
		}
	}
	return nil, ErrNotFound
}

func (db *Database) SaveAccount(acc *models.Account) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Accounts[acc.ID] = acc
	return db.Save()
}
