package mongo

import (
	"github.com/go-mixins/mongodb"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/andviro/filer/backend"
)

// Errors further specifies backend error class
var Errors = backend.Errors.Sub("mongo")

// Backend implements data store in MongoDB
type Backend struct {
	db *mongodb.DB
}

var _ backend.Backend = (*Backend)(nil)

// New returns new backend instance
func New(db *mongodb.DB) (res *Backend, err error) {
	if err = db.Session.Ping(); err != nil {
		err = Errors.Wrap(err, "connecting")
		return
	}
	res = &Backend{db.Clone()}
	return
}

// Close closes DB connection
func (b *Backend) Close() error {
	b.db.Close()
	return nil
}

// Stat requests file metadata from store
func (b *Backend) Stat(filename string) (res *backend.FileInfo, err error) {
	db := b.db.Clone()
	defer db.Close()
	coll := db.C("meta")
	var fi FileInfo
	res = &fi.FileInfo
	err = coll.Find(bson.M{
		"names": filename,
		"state": "saved",
	}).One(&fi)
	if err == mgo.ErrNotFound {
		err = backend.ErrNotFound.Errorf("%q not found", filename)
		return
	}
	if err != nil {
		err = Errors.Wrapf(err, "loading meta for %q", filename)
		return
	}
	return
}

// CreateTransaction executes commit function in Create transaction context
func (b *Backend) CreateTransaction(filename string, commit func(*backend.FileInfo) error) (err error) {
	db := b.db.Clone()
	defer db.Close()
	coll := db.C("meta")
	fi := &FileInfo{
		ID: bson.NewObjectId(),
		FileInfo: backend.FileInfo{
			Names:        []string{filename},
			LastModified: bson.Now(),
		},
		State: "writing",
	}
	err = coll.Insert(fi)
	if mgo.IsDup(err) {
		err = backend.ErrBusy.Errorf("%q is busy", filename)
		return
	}
	if err != nil {
		err = Errors.Wrapf(err, "inserting %q in db", filename)
		return
	}
	if err = commit(&fi.FileInfo); err != nil {
		if e := coll.RemoveId(fi.ID); e != nil {
			return Errors.Wrapf(e, "deleting %q (original error: %v)", filename, err)
		}
		return
	}
	return Errors.Wrapf(coll.UpdateId(fi.ID, bson.M{
		"$set": bson.M{
			"state":    "saved",
			"disksize": fi.DiskSize,
			"fileid":   fi.FileID,
		},
	}), "updating record for %q", filename)
}

// RemoveTransaction executes commit function in Remove transaction context
func (b *Backend) RemoveTransaction(filename string, commit func(string, *backend.FileInfo) error) (err error) {
	db := b.db.Clone()
	defer db.Close()
	coll := db.C("meta")
	var fi FileInfo
	_, err = coll.Find(bson.M{
		"names": filename,
		"state": "saved",
	}).Apply(mgo.Change{
		Update: bson.M{
			"$pull": bson.M{
				"names": filename,
			},
			"$set": bson.M{
				"state": "removing",
			},
		},
		ReturnNew: true,
	}, &fi)
	if err == mgo.ErrNotFound {
		err = backend.ErrNotFound.Errorf("%q not found", filename)
		return
	}
	if err != nil {
		err = Errors.Wrapf(err, "updating %q", filename)
		return
	}
	if err = commit(filename, &fi.FileInfo); err != nil {
		if e := coll.UpdateId(fi.ID, bson.M{
			"$addToSet": bson.M{
				"names": filename,
			},
			"$set": bson.M{
				"state": "saved",
			},
		}); e != nil {
			return Errors.Wrapf(e, "updating %q (original error: %v)", filename, err)
		}
		return
	}
	if len(fi.Names) == 0 {
		return Errors.Wrapf(coll.RemoveId(fi.ID), "deleting %q", filename)
	}
	return Errors.Wrapf(coll.UpdateId(fi.ID, bson.M{
		"$set": bson.M{
			"state": "saved",
		},
	}), "updating state for %q", filename)
}
