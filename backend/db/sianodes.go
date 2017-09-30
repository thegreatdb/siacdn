package db

import (
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/backend/models"
)

func (db *Database) GetSiaNode(id uuid.UUID) (*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if sn, ok := db.SiaNodes[id]; ok {
		return sn.Copy(), nil
	} else {
		return nil, ErrNotFound
	}
}

func (db *Database) GetSiaNodeByShortcode(shortcode string) (*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sn := range db.SiaNodes {
		if sn.Shortcode == shortcode {
			return sn.Copy(), nil
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

func (db *Database) GetSiaNodesByAccount(id uuid.UUID) ([]*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	nodes := []*models.SiaNode{}
	for _, sn := range db.SiaNodes {
		if sn.AccountID == id {
			nodes = append(nodes, sn.Copy())
		}
	}
	return nodes, nil
}

func (db *Database) GetOrphanedSiaNode(accountID uuid.UUID) (*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sn := range db.SiaNodes {
		if sn.AccountID == accountID &&
			sn.MinioInstancesActivated == 0 ||
			sn.MinioInstancesActivated < sn.MinioInstancesRequested {
			return sn.Copy(), nil
		}
	}
	return nil, ErrNotFound
}

func (db *Database) GetPendingSiaNodes() ([]*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	nodes := []*models.SiaNode{}
	for _, sn := range db.SiaNodes {
		if sn.Pending() {
			nodes = append(nodes, sn.Copy())
		}
	}
	return nodes, nil
}

func (db *Database) GetReadyOrphanedSiaNodes() ([]*models.SiaNode, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	nodes := []*models.SiaNode{}
	for _, sn := range db.SiaNodes {
		if sn.Status == models.SIANODE_STATUS_READY &&
			sn.MinioInstancesActivated == 0 ||
			sn.MinioInstancesActivated < sn.MinioInstancesRequested {
			nodes = append(nodes, sn.Copy())
		}
	}
	return nodes, nil
}
