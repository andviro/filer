package filer

import (
	"io"
	"time"
)

// FileInfo is a subset of os.FileInfo
type FileInfo interface {
	Size() int64        // Size in bytes
	ModTime() time.Time // Last modified time
}

// Storage enumerate basic file operations available in all storage backends
type Storage interface {
	Open(fn string) (res io.ReadCloser, err error)
	Create(fn string) (res io.WriteCloser, err error)
	Remove(fn string) error
	Rename(from, to string) error
}

// Meta provides file metadata services
type Meta interface {
	Stat(fn string) (res FileInfo, err error)
}

//go:generate moq -out mock/filer.go -pkg mock . FileInfo Storage Backend
