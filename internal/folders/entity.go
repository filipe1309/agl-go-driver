package folders

import (
	"database/sql"
	"errors"
	"time"

	"github.com/filipe1309/agl-go-driver/internal/common"
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
		folder.ParentID = common.NullInt64{sql.NullInt64{Int64: parent_id, Valid: true}}
	}

	err := folder.Validate()
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

type Folder struct {
	ID        int64     `json:"id"`
	ParentID  common.NullInt64 `json:"parent_id"`
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
