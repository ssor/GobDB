// Implements a persistant key-value store of gob-compatible
// types. This is accomplished with a light wrapper around 
// leveldb and Go's gob encoding library.
package GobDB

import (
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)


// LevelDB wrapper.
type DB struct {
	internal *leveldb.DB
	location string
	encoder FilteredEncoder
	decoder Decoder
	prepared bool
}


// Returns unopened database at with given datafile.
func At(path string) *DB {
	return &DB {location: path}
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
func (db *DB) Put(key, value interface{}) error {
	// Encode key via gob, registering types if necessary.
	o1, err := db.encode(key)
	if err != nil {
		return err
	}

	// Encode value via gob, registering types if necessary.
	o2, err := db.encode(value)
	if err != nil {
		return err
	}

	// Insert gobbed values into leveldb.
	return db.internal.Put(o1, o2, nil)
}


// // Encodes given key via gob, fetches the corresponding value from
// // within leveldb, and decodes that value into parameter two.
func (db *DB) Get(key, value interface{}) error {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return err
	}

	// Fetch gob-encoded values.
	val, err := db.internal.Get(obj, nil)
	if err != nil {
		return err
	}

	// Decode value into second paramater (which should be a pointer)
	return db.decoder.Decode(val, value)
}


// // Encodes given key via gob and checks if the resulting byte 
// // slice exists in the database's internal leveldb.
func (db DB) Contains(key interface{}) bool {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return false
	}

	// Note: this is niave - should check error type.
	_, err = db.internal.Get(obj, nil)
	return err == nil
}


// // Encodes given key via gob, deleting the resulting byte slice 
// // from the database's internal leveldb.
func (db *DB) Delete(key interface{}) error {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return err
	}

	// Delete!
	return db.internal.Delete(obj, nil)
}


// Encodes given key via gob, registers its type if necessary,
// and routes any errors outward.
func (db *DB) encode(key interface{}) ([]byte, error) {
	err := db.Open()
	if err != nil {
		return []byte{}, err
	}

	// Gob encode key.
	def, obj, err := db.encoder.Encode(key)
	if err != nil {
		return []byte{}, err
	}

	// Register key type.
	err = db.registerType(def, obj)
	if err != nil {
		return []byte{}, err
	}

	return obj, nil
}



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



func (db *DB) typeIndex(def, obj []byte) int {
	key := []byte("type:")
	key = append(key, def...)
	key = append(key, obj...)
	value, err := db.internal.Get(key, nil)
	if err != nil {
		return -1
	}
	
	n, err := strconv.Atoi(string(value))
	if err != nil {
		return -1
	}
	return n
}



func (db *DB) registerType(def, obj []byte) error {
	// i := db.typeIndex(def, obj)
	// if i == -1 {
	// 	val := []byte(strconv.Itoa(db.prepSize()))
	// 	key := typeKey(def, obj)
	// 	err := db.internal.Put(key, val)
	// 	if err != nil {
	// 		return err
	// 	} 
	// 	db.incPrepSize()
	// }
	// return nil
	return nil
}


