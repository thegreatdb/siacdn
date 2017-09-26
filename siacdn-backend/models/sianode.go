package models

import (
	"time"

	"github.com/google/uuid"
)

const SIANODE_STATUS_CREATED = "created"           // This object has been created
const SIANODE_STATUS_DEPLOYED = "deployed"         // Deployment yaml has been sent to kube
const SIANODE_STATUS_INSTANCED = "instanced"       // Pod instance has been seen
const SIANODE_STATUS_SNAPSHOTTED = "snapshotted"   // Bootstrap snapshot script has finished
const SIANODE_STATUS_SYNCHRONIZED = "synchronized" // Blockchain has synchronized
const SIANODE_STATUS_INITIALIZED = "initialized"   // Wallet has been initialized
const SIANODE_STATUS_FUNDED = "funded"             // Account has received initial funding
const SIANODE_STATUS_CONFIGURED = "configured"     // Allowance has been set
const SIANODE_STATUS_READY = "ready"               // Everything is ready to go and contracts are all set
const SIANODE_STATUS_STOPPED = "stopped"           // A user or administrator has stopped the node
const SIANODE_STATUS_DEPLETED = "depleted"         // All the SiaCoins in the accoint have been used up
const SIANODE_STATUS_ERROR = "error"               // An error has occurred in the process

type SiaNode struct {
	ID          uuid.UUID `json:"id"`
	Shortcode   string    `json:"shortcode"`
	Status      string    `json:"status"`
	CreatedTime time.Time `json:"created_time"`
}

func NewSiaNode() (*SiaNode, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	shortcode := id.String() // TODO: This needs to make it shorter lol
	return &SiaNode{
		ID:          id,
		Shortcode:   shortcode,
		Status:      SIANODE_STATUS_CREATED,
		CreatedTime: time.Now().UTC(),
	}, nil
}
