package disk

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andviro/filer/storage"
)

// Errors represent disk errors sub-class
var Errors = storage.Errors.Sub("disk")

// Storage provides on-disk file storage primitives
type Storage struct {
	basePath string
}

var _ storage.Storage = (*Storage)(nil)

// New creates disk storage with specified base path
func New(basePath string) *Storage {
	return &Storage{basePath}
}

// Save file from source reader
func (s *Storage) Save(fn string, src io.Reader) (n int64, id string, err error) {
	dir, base := filepath.Split(fn)
	dir = filepath.Join(s.basePath, strings.TrimPrefix(dir, s.basePath))
	if err = os.MkdirAll(dir, 0755); err != nil {
		err = Errors.Wrap(err, "ensuring path")
		return
	}
	id = filepath.Join(dir, base)
	dest, err := os.Create(filepath.Join(dir, base))
	if err != nil {
		err = Errors.Wrap(err, "creating file")
		return
	}
	n, err = io.Copy(dest, src)
	err = Errors.Wrap(err, "saving file")
	return
}

// Load file into dest writer
func (s *Storage) Load(id string, dest io.Writer) (err error) {
	src, err := os.Open(id)
	if err != nil {
		err = Errors.Wrap(err, "opening file")
		return
	}
	_, err = io.Copy(dest, src)
	err = Errors.Wrap(err, "loading file")
	return
}

// Remove file
func (s *Storage) Remove(id string) error {
	err := os.Remove(id)
	if err == os.ErrNotExist {
		return storage.ErrNotFound.Wrapf(err, "removing %q", id)
	}
	return Errors.Wrapf(err, "removing %q", id)
}
