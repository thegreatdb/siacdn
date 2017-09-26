package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/siacdn-backend/models"
)

var ErrNotFound = errors.New("Database object not found")

type Database struct {
	Accounts   map[uuid.UUID]*models.Account   `json:"accounts"`
	AuthTokens map[uuid.UUID]*models.AuthToken `json:"auth_tokens"`
	SiaNodes   map[uuid.UUID]*models.SiaNode   `json:"sianodes"`

	filepath string
	mu       *sync.RWMutex
}

func (db *Database) EnsureDefaults(filepath string) {
	if db.Accounts == nil {
		db.Accounts = map[uuid.UUID]*models.Account{}
	}
	if db.AuthTokens == nil {
		db.AuthTokens = map[uuid.UUID]*models.AuthToken{}
	}
	if db.SiaNodes == nil {
		db.SiaNodes = map[uuid.UUID]*models.SiaNode{}
	}
	if filepath != "" && db.filepath != filepath {
		db.filepath = filepath
	}
	db.mu = &sync.RWMutex{}
}

func (db *Database) Save() error {
	// Ensure the file directory exists
	dir := filepath.Dir(db.filepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0644); err != nil {
			return err
		}
	}

	// Serialize the database to json
	data, err := json.Marshal(db)
	if err != nil {
		return err
	}

	// Write the json to disk
	return ioutil.WriteFile(db.filepath, data, 0644)
}

func OpenDatabase(filepath string) (*Database, error) {
	db := &Database{}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("Could not load file (" + err.Error() + ")")
		log.Println("Continuing with empty database...")
		db.EnsureDefaults(filepath)
		return db, nil
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		log.Println("Could not decode file, continuing with empty database: " + err.Error())
		// Minor realloc in case some weird implementation of unmarshal is destructive
		db = &Database{}
	}
	db.EnsureDefaults(filepath)
	return db, nil
}
