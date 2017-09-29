package db

import (
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

func (db *Database) GetWalletSeed(siaNodeID uuid.UUID) (*models.WalletSeed, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if seed, ok := db.WalletSeeds[siaNodeID]; ok {
		return seed.Copy(), nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) SaveWalletSeed(seed *models.WalletSeed) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.WalletSeeds[seed.SiaNodeID] = seed
	return db.Save()
}
