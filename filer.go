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

// Service provides all file-related primitives
type Service interface {
	Download(fn string, dest io.Writer) (err error)
	Upload(fn string, src io.Reader) (err error)
	Remove(fn string) error
	Rename(from, to string) error
	Stat(fn string) (res FileInfo, err error)
}

//go:generate moq -out mock/filer.go -pkg mock . FileInfo Storage Backend
