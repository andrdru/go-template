package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Session struct {
	ID        int64        `json:"-"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt *time.Time   `json:"-"`
	UserID    int64        `json:"user_id"`
	Token     string       `json:"token"`
	Extra     SessionExtra `json:"-"`

	User  *User  `json:"-"`
	Email string `json:"-"`
	Pass  string `json:"-"`
}

type SessionExtra struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

var (
	// ErrUnknownFieldDataType .
	ErrUnknownFieldDataType = errors.New("unknown field data type")
)

// Scan implement sql.Scanner
func (l *SessionExtra) Scan(src interface{}) (err error) {
	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	default:
		return ErrUnknownFieldDataType
	}

	err = json.Unmarshal(source, &l)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

// Value implement sql/driver.Valuer
func (l SessionExtra) Value() (driver.Value, error) {
	return json.Marshal(l)
}
