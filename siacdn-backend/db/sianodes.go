package db

import (
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

func (db *Database) GetSiaNode(id uuid.UUID) (*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if sn, ok := db.SiaNodes[id]; ok {
		return sn, nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) GetSiaNodeByShortcode(shortcode string) (*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sn := range db.SiaNodes {
		if sn.Shortcode == shortcode {
			return sn, nil
		}
	}
	return nil, ErrNotFound
}

func (db *Database) SaveSiaNode(sn *models.SiaNode) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.SiaNodes[sn.ID] = sn
	return db.Save()
}
