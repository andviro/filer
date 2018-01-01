package disk

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andviro/filer/backend"
)

// Backend takes all file metadata from OS
type Backend struct {
	basePath string
}

var _ backend.Backend = (*Backend)(nil)

// Errors sub-class for disk backend
var Errors = backend.Errors.Sub("disk")

// New creates disk backend with specified base path
func New(basePath string) *Backend {
	return &Backend{basePath}
}

// Stat returns file info
func (be *Backend) Stat(filename string) (res *backend.FileInfo, err error) {
	id := filepath.Join(be.basePath, strings.TrimPrefix(filename, be.basePath))
	st, err := os.Stat(id)
	if err == os.ErrNotExist {
		err = backend.ErrNotFound.Wrapf(err, "stat %q", id)
		return
	}
	if err != nil {
		err = Errors.Wrapf(err, "stat %q", id)
		return
	}
	res = &backend.FileInfo{
		DiskSize:     st.Size(),
		Names:        []string{filename},
		LastModified: st.ModTime(),
		FileID:       id,
	}
	return
}

// CreateTransaction calls commit on specific fileinfo from disk
func (be *Backend) CreateTransaction(filename string, commit func(*backend.FileInfo) error) (err error) {
	fi := &backend.FileInfo{
		Names:        []string{filename},
		LastModified: time.Now(),
	}
	return commit(fi)
}

// RemoveTransaction calls commit on specific fileinfo from disk
func (be *Backend) RemoveTransaction(filename string, commit func(*backend.FileInfo) error) (err error) {
	fi, err := be.Stat(filename)
	if err != nil {
		return
	}
	return commit(fi)
}

// Rename renames file on disk
func (be *Backend) Rename(from, to string) (err error) {
	from = filepath.Join(be.basePath, strings.TrimPrefix(from, be.basePath))
	to = filepath.Join(be.basePath, strings.TrimPrefix(to, be.basePath))
	if err = os.Rename(from, to); err == os.ErrNotExist {
		return backend.ErrNotFound.Wrapf(err, "renaming %q to %q", from, to)
	}
	return Errors.Wrapf(err, "renaming %q to %q", from, to)
}
