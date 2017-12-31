package disk

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andviro/filer"
	"github.com/andviro/filer/storage"
)

// Errors represent disk errors sub-class
var Errors = storage.Errors.Sub("disk")

// Storage provides on-disk file storage primitives
type Storage struct {
	DataDir string `default:"/tmp"` // File base path
}

var _ storage.Storage = (*Storage)(nil)

// Ensure concatenates path components with service path prefix and creates
// directory on disk. Final combined path is returned.
func (s *Storage) Ensure(path ...string) (res string, err error) {
	res = strings.TrimPrefix(filepath.Join(path...), s.DataDir)
	res = filepath.Join(s.DataDir, res)
	err = Errors.Wrap(os.MkdirAll(res, 0755), "ensuring path")
	return
}

// Create file for writing
func (s *Storage) Create(path string) (res io.WriteCloser, err error) {
	res, err = os.Create(path)
	err = Errors.Wrap(err, "creating file")
	return
}

// Open file for reading
func (s *Storage) Open(fn string) (res io.ReadCloser, err error) {
	res, err = os.Open(fn)
	err = Errors.Wrap(err, "opening file")
	return
}

// Stat returns file info from disk
func (s *Storage) Stat(fn string) (res filer.FileInfo, err error) {
	res, err = os.Stat(fn)
	err = Errors.Wrap(err, "opening file")
	return
}

// Remove file from disk
func (s *Storage) Remove(fn string) error {
	return Errors.Wrap(os.RemoveAll(fn), "removing")
}

// Rename file
func (s *Storage) Rename(from, to string) error {
	return Errors.Wrap(os.Rename(from, to), "renaming")
}
