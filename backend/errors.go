package backend

import "github.com/andviro/filer"

// Errors specifies error class for a backend
var Errors = filer.Errors.Sub("backend")

// ErrNotFound class is returned on file not found
var ErrNotFound = Errors.Sub("not found")

// ErrBusy class is returned on file access collision
var ErrBusy = Errors.Sub("busy")
