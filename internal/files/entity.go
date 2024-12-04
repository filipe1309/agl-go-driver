package files

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v4"
)

var (
	ErrNameEmpty    = errors.New("name is required")
	ErrPathEmpty    = errors.New("path is required")
	ErrOwnerIDEmpty = errors.New("owner_id is required")
	ErrTypeEmpty    = errors.New("type is required")
)

func New(owner_id int64, name, fileType, path string) (*File, error) {
	file := File{
		OwnerID:   owner_id,
		Name:      name,
		Type:      fileType,
		Path:      path,
		UpdatedAt: time.Now(),
	}

	err := file.Validate()
	if err != nil {
		return nil, err
	}

	return &file, nil
}

type File struct {
	ID        int64     `json:"id"`
	FolderID  null.Int  `json:"-"`
	OwnerID   int64     `json:"owner_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Path      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"-"`
}

func (f *File) Validate() error {
	if f.Name == "" {
		return ErrNameEmpty
	}

	if f.Path == "" {
		return ErrPathEmpty
	}

	if f.OwnerID == 0 {
		return ErrOwnerIDEmpty
	}

	if f.Type == "" {
		return ErrTypeEmpty
	}

	return nil
}
