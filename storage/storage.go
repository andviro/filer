package storage

import "io"

// Storage provides file storage operations
type Storage interface {
	Save(fn string, src io.Reader) (n int64, id string, err error)
	Load(id string, dest io.Writer) (err error)
	Remove(id string) error
}
