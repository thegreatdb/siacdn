package db

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

func (db *Database) GetAccount(id uuid.UUID) (*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if account, ok := db.Accounts[id]; ok {
		return account, nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) GetAccountByUsername(username string) (*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, acc := range db.Accounts {
		if strings.ToLower(acc.Username) == strings.ToLower(username) {
			return acc, nil
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
