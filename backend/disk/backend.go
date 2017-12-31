package disk

import (
	"os"
	"time"

	"github.com/andviro/filer/backend"
)

// Backend takes all file metadata from OS
type Backend struct{}

var _ backend.Backend = (*Backend)(nil)

// Errors sub-class for disk backend
var Errors = backend.Errors.Sub("disk")

// Stat returns file info
func (be *Backend) Stat(filename string) (res *backend.FileInfo, err error) {
	st, err := os.Stat(filename)
	if err == os.ErrNotExist {
		err = backend.ErrNotFound.Wrapf(err, "stat %q", filename)
		return
	}
	if err != nil {
		err = Errors.Wrapf(err, "stat %q", filename)
		return
	}
	res = &backend.FileInfo{
		DiskSize:     st.Size(),
		Names:        []string{filename},
		LastModified: st.ModTime(),
		FileID:       "file://" + filename,
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
func (be *Backend) RemoveTransaction(filename string, commit func(string, *backend.FileInfo) error) (err error) {
	fi, err := be.Stat(filename)
	if err != nil {
		return
	}
	return commit(filename, fi)
}