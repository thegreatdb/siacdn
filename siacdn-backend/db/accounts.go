package db

import (
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

func (db *Database) GetAccount(id uuid.UUID) (*models.Account, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if account, ok := db.Accounts[id]; ok {
		return account, nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) SaveAccount(acc *models.Account) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Accounts[acc.ID] = acc
	return db.Save()
}
