package buckettest

import (
	"io"
	"os"
)

type MockBucket struct {
	content map[string][]byte
}

func New() *MockBucket {
	return &MockBucket{
		content: make(map[string][]byte),
	}
}

func (mb *MockBucket) Upload(file io.Reader, key string) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	mb.content[key] = data

	return nil
}

func (mb *MockBucket) Download(src, dest string) (*os.File, error) {
	return nil, nil
}

func (mb *MockBucket) Delete(key string) error {
	return nil
}
