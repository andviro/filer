package storage

import "github.com/andviro/filer"

// Errors represent storage errors sub-class
var Errors = filer.Errors.Sub("storage")

// ErrNotFound is returned by storage if file not found
var ErrNotFound = Errors.Sub("not found")
