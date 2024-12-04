package bucket

import (
	"io"
)

type MockBucket struct {
	content map[string][]byte
}

func (mb *MockBucket) Upload(file io.Reader, key string) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	mb.content[key] = data

	return nil
}

func (mb *MockBucket) Download(src, dest string) error {
	return nil
}

func (mb *MockBucket) Delete(key string) error {
	return nil
}
