package models

import (
	"time"

	"github.com/google/uuid"
)

type WalletSeed struct {
	SiaNodeID   uuid.UUID `json:"sianode_id"`
	Words       string    `json:"words"`
	CreatedTime time.Time `json:"created_time"`
}

func NewWalletSeed(siaNodeId uuid.UUID, words string) (*WalletSeed, error) {
	return &WalletSeed{
		SiaNodeID:   siaNodeId,
		Words:       words,
		CreatedTime: time.Now().UTC(),
	}, nil
}

func (seed *WalletSeed) Copy() *WalletSeed {
	var cpy WalletSeed
	cpy = *seed
	return &cpy
}
