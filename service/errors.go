package service

import "github.com/andviro/filer"

// Errors represent service-specific error class
var Errors = filer.Errors.Sub("service")

// ErrNotFound is returned by service if file not found
var ErrNotFound = Errors.Sub("not found")
