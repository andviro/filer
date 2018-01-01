package service

import (
	"io"

	"github.com/andviro/filer"
	"github.com/andviro/filer/backend"
	"github.com/andviro/filer/storage"
)

// Service implements filer.Service interface
type Service struct {
	backend backend.Backend
	storage storage.Storage
}

var _ filer.Service = (*Service)(nil)

// New creates Service instance with provided storages and backends
func New(b backend.Backend, s storage.Storage) *Service {
	return &Service{b, s}
}

// Stat returns file metadata, ErrNotFound if not found
func (s *Service) Stat(fn string) (res filer.FileInfo, err error) {
	return s.backend.Stat(fn)
}

// Upload creates file entry with specified filename from source reader
func (s *Service) Upload(fn string, src io.Reader) error {
	return s.backend.CreateTransaction(fn, func(fi *backend.FileInfo) (err error) {
		fi.DiskSize, fi.FileID, err = s.storage.Save(fn, src)
		return
	})
}

// Download copies file with specified filename into dest writer
func (s *Service) Download(fn string, dest io.Writer) (err error) {
	fi, err := s.backend.Stat(fn)
	if err != nil {
		return
	}
	return s.storage.Load(fi.FileID, dest)
}

// Remove deletes file from storage ignoring file not found errors
func (s *Service) Remove(fn string) (err error) {
	err = s.backend.RemoveTransaction(fn, func(fi *backend.FileInfo) (err error) {
		if err = s.storage.Remove(fi.FileID); storage.ErrNotFound.Contains(err) {
			return nil
		}
		return
	})
	if backend.ErrNotFound.Contains(err) {
		err = nil
	}
	return
}

// Rename changes name of file
func (s *Service) Rename(from, to string) error {
	return nil
}
