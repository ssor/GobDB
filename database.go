// Implements a persistant key-value store of gob-compatible
// types. This is accomplished with a light wrapper around 
// leveldb and Go's gob encoding library.
//
// NOTE: this library is not yet goroutine-safe.
package GobDB

import (
	"bytes"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)


// LevelDB wrapper.
type DB struct {
	internal *leveldb.DB
	location string
	essentials *bytes.Buffer
	encoder FilteredEncoder
	prepared bool
}


// Returns unopened database at with given datafile.
func At(path string) DB {
	buffer := bytes.NewBuffer([]byte{})
	return DB {
		location: path,
		essentials: buffer,
		encoder: MakeFilteredEncoder(buffer),
	}
}


// Opens database if not done already.
func (db *DB) Open() error {
	if db.IsOpen() {
		return nil
	}

	ret, err := leveldb.OpenFile(db.location, nil)
	if err == nil {
		db.internal = ret
		db.prepareEncoder()
	}
	return err
}


// 
func (db DB) IsOpen() bool {
	return db.internal != nil
}


// Closes database if not done already.
func (db *DB) Close() {
	if db.IsOpen() {
		db.internal.Close()
		db.internal = nil
	}
}


// Encodes given key and value through gob, inserting resulting
// byte slices into the database's internal leveldb.
// func (db *DB) Put(key, value interface{}) error {
// 	err := db.Open()
// 	if err != nil {
// 		return err
// 	}


// }


// // Encodes given key via gob, fetches the corresponding value from
// // within leveldb, and decodes that value into parameter two.
// func (db *DB) Get(key, value interface{}) error {
// 	err := db.Open()
// 	if err != nil {
// 		return err
// 	}

// 	gk, _ := db.encoder.Encode(key)
// 	gv, err := 	
// }


// // Encodes given key via gob and checks if the resulting byte 
// // slice exists in the database's internal leveldb.
// func (db DB) Contains(key interface{}) bool {

// }


// // Encodes given key via gob, deleting the resulting byte slice 
// // from the database's internal leveldb.
// func (db *DB) Delete(key interface{}) error {

// }



// When the database is opened for the first time, scrolls through
// all entries to form the same encoder that was used before the 
// previous db.Close() call.
//
// Note that this is an expensive operation, scaling linearly with
// dataset size. Accordingly, you should utilize the open and close
// methods instead of initializing new DB objects all the time. For 
// each initialization, the decoder is literally thrown in the 
// garbage.
func (db *DB) prepareEncoder() error {
	if db.prepared == true {
		return nil
	}


	iter := db.internal.NewIterator(&util.Range{Start: []byte("prep:"), Limit: []byte("prep::")}, nil)
	for iter.Next() {
		value := iter.Value()
		db.encoder.Encode(value)
	}
	iter.Release()
	err := iter.Error()
	if err == nil {
		db.prepared = true
	}
	return err
}


func (db *DB) setPrepSize(value int) error {
	key := []byte("nprep")
	data := []byte(strconv.Itoa(value))
	return db.internal.Put(key, data, nil)
}


func (db *DB) incPrepSize() error {
	size := db.prepSize()
	if size == -1 {
		return db.setPrepSize(1)
	} else {
		return db.setPrepSize(size + 1)
	}
}


func (db *DB) prepSize() int {
	err := db.Open()
	if err != nil {
		return -1
	}

	val, err := db.internal.Get([]byte("nprep"), nil)
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(string(val))
	return n
}

