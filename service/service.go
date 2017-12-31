package service

import (
	"io"

	"github.com/andviro/filer"
	"github.com/andviro/filer/backend"
	"github.com/andviro/filer/storage"
)

// Service implements filer.Service interface
type Service struct {
	backends []backend.Backend
	storages map[string]storage.Storage
}

var _ filer.Service = (*Service)(nil)

// New creates Service instance with provided storages and backends
func New(b []backend.Backend, s map[string]storage.Storage) *Service {
	return &Service{b, s}
}

// Stat returns file metadata, ErrNotFound
func (s *Service) Stat(fn string) (res filer.FileInfo, err error) {
	for _, b := range s.backends {
		res, err = b.Stat(fn)
		if backend.ErrNotFound.Contains(err) {
			continue
		}
		if err != nil {
			return
		}
		break
	}
	return
}

// Create creates file entry in backend and storage
func (s *Service) Create(fn string) (res io.WriteCloser, err error) {
	return
}

// Open opens file for reading
func (s *Service) Open(fn string) (res io.ReadCloser, err error) {
	return
}

// Remove deletes file from storage
func (s *Service) Remove(fn string) error {
	return nil
}

// Rename changes name of file
func (s *Service) Rename(from, to string) error {
	return nil
}
