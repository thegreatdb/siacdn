package db

import (
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

func (db *Database) GetAuthToken(id uuid.UUID) (*models.AuthToken, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if tok, ok := db.AuthTokens[id]; ok {
		return tok.Copy(), nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) SaveAuthToken(tok *models.AuthToken) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.AuthTokens[tok.ID] = tok
	return db.Save()
}
