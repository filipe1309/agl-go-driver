package folders

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNameEmpty = errors.New("name is required")
)

func New(name string, parent_id int64) (*Folder, error) {
	folder := Folder{
		Name:      name,
		UpdatedAt: time.Now(),
	}

	if parent_id > 0 {
		folder.ParentID = NullInt64{sql.NullInt64{Int64: parent_id, Valid: true}}
	} 

	err := folder.Validate()
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

type Folder struct {
	ID        int64     `json:"id"`
	ParentID  NullInt64 `json:"parent_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"-"`
}

func (f *Folder) Validate() error {
	if f.Name == "" {
		return ErrNameEmpty
	}

	return nil
}

type FolderContent struct {
	Folder  Folder           `json:"folder"`
	Content []FolderResource `json:"content"`
}

type FolderResource struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
