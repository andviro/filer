package backend

// Backend provides storage of file metadata
type Backend interface {
	Stat(filename string) (res *FileInfo, err error)
	CreateTransaction(filename string, commit func(*FileInfo) error) (err error)
	RemoveTransaction(filename string, commit func(string, *FileInfo) error) (err error)
}
