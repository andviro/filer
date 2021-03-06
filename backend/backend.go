package backend

// Backend provides storage of file metadata
type Backend interface {
	Stat(filename string) (res *FileInfo, err error)
	CreateTransaction(filename string, commit func(*FileInfo) error) (err error)
	RemoveTransaction(filename string, commit func(*FileInfo) error) (err error)
	Rename(from, to string) (err error)
}

//go:generate moq -out mock/backend.go -pkg mock . Backend
