package backend

import (
	"time"

	"github.com/andviro/filer"
)

// FileInfo contains essential file metadata
type FileInfo struct {
	DiskSize     int64     // Physical size
	Names        []string  // Filename aliases
	LastModified time.Time // Creation/modification time
	FileID       string    // Opaque storage identifier
}

var _ filer.FileInfo = (*FileInfo)(nil)

// Size returns file size in bytes
func (fi *FileInfo) Size() int64 {
	return fi.DiskSize
}

// ModTime returns last update time for a file
func (fi *FileInfo) ModTime() time.Time {
	return fi.LastModified
}
