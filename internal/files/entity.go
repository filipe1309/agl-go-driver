package files

import (
	"errors"
	"time"
)

var (
	ErrNameEmpty = errors.New("name is required")
	ErrPathEmpty = errors.New("path is required")
	ErrOwnerIDEmpty = errors.New("owner_id is required")
	ErrTypeEmpty = errors.New("type is required")
)

func New(name, path string, folder_id, owner_id int64) (*File, error) {
	file := File{
		Name:     name,
		FolderID: folder_id,
		OwnerID:  owner_id,
	}

	err := file.Validate()
	if err != nil {
		return nil, err
	}

	return &file, nil
}

type File struct {
	ID        int64     `json:"id"`
	FolderID  int64     `json:"-"`
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
