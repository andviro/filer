package mongo

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/andviro/filer/backend"
)

// FileInfo stores metadata specifically for Mongo backend
type FileInfo struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	State            string
	backend.FileInfo `bson:",inline"`
}
