package client

import (
	"github.com/andviro/filer"
)

// Client implements filer.Service over HTTP connection
type Client struct {
	filer.Storage
}
