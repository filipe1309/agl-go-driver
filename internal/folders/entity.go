package folders

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v4"
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
		folder.ParentID = null.IntFrom(parent_id)
		// common.NullInt64{Int64: parent_id, Valid: true}
	}

	err := folder.Validate()
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

type Folder struct {
	ID int64 `json:"id"`
	// ParentID  common.NullInt64 `json:"parent_id"`
	ParentID  null.Int  `json:"parent_id"`
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
