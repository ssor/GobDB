// Package GobDB implements a persistant key-value store of
// gob-compatible types. This is accomplished with a light
// wrapper around leveldb and Go's gob encoding library.
package GobDB

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strconv"
)

// DB is a LevelDB wrapper that stores key-value pairs of
// gob-compatible types.
type DB struct {
	internal *leveldb.DB
	location string
	encoder  FilteredEncoder
	decoder  Decoder
	prepared bool
}

// At returns an unopened database at with given datafile.
func At(path string) *DB {
	return &DB{location: path}
}

// Open sets up the internal leveldb if not done already.
func (db *DB) Open() error {
	if db.IsOpen() {
		return nil
	}

	ret, err := leveldb.OpenFile(db.location, nil)
	if err == nil {
		db.internal = ret
		db.prepare()
	}
	return err
}

// IsOpen checks whether or not the database is open.
func (db DB) IsOpen() bool {
	return db.internal != nil
}

// Close tears down the internal leveldb, writing all contents
// to file.
func (db *DB) Close() {
	if db.IsOpen() {
		db.internal.Close()
		db.internal = nil
	}
}

// Put encodes given key and value through gob, inserting resulting
// byte slices into the database's internal leveldb.
func (db *DB) Put(key, value interface{}) error {
	// Encode key via gob, registering types if necessary.
	o1, err := db.encode(key)
	if err != nil {
		return err
	}

	// Form prefixed key.
	pkey := []byte("key:")
	pkey = append(pkey, o1...)

	// Encode value via gob, registering types if necessary.
	val, err := db.encode(value)
	if err != nil {
		return err
	}

	// Insert gobbed values into leveldb.
	return db.internal.Put(pkey, val, nil)
}

// Get ncodes given key via gob, fetches the corresponding
// value from within leveldb, and decodes that value into
// parameter two.
func (db *DB) Get(key, value interface{}) error {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return err
	}

	// Form prefixed key.
	pkey := []byte("key:")
	pkey = append(pkey, obj...)

	// Fetch gob-encoded value.
	val, err := db.internal.Get(pkey, nil)
	if err != nil {
		return err
	}

	// Decode value into second paramater (which should be a
	// pointer)
	return db.decoder.Decode(val, value)
}

// Has encodes given key via gob and checks if the resulting
// byte slice exists in the database's internal leveldb.
func (db DB) Has(key interface{}) bool {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return false
	}

	// Note: this is niave - should check error type.
	_, err = db.internal.Get(obj, nil)
	return err == nil
}

// Delete encodes given key via gob, deleting the resulting
// byte slice from the database's internal leveldb.
func (db *DB) Delete(key interface{}) error {
	// Encode key via gob, registering its type if necessary.
	obj, err := db.encode(key)
	if err != nil {
		return err
	}

	// Form key bytes.
	kbytes := []byte("key:")
	kbytes = append(kbytes, obj...)

	// Delete!
	return db.internal.Delete(kbytes, nil)
}

// Entries counts key-value pairs in the database.
func (db *DB) Entries() int {
	i := 0
	iter := db.internal.NewIterator(util.BytesPrefix([]byte("key:")), nil)
	for iter.Next() {
		i++
	}
	iter.Release()
	return i
}

// Reset erases caches and closes leveldb. This way, the db
// is forced to reload gobbed values as though it had just 
// been opened for the first time.
func (db *DB) Reset() {
	db.Close()
	db.prepared = false
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

	// Register key with decoder.
	db.decoder.Register(append(def, obj...))

	// Register key type.
	err = db.registerType(def, obj)
	if err != nil {
		return []byte{}, err
	}

	// err = db.decoder.Decode(append(def, obj...), nil)
	// if err != nil {
	// 	return []byte{}, err
	// }

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
func (db *DB) prepare() error {
	if db.prepared == true {
		return nil
	}

	iter := db.internal.NewIterator(util.BytesPrefix([]byte("prep#")), nil)
	for iter.Next() {
		value := iter.Value()
		db.encoder.Encode(value)
		db.decoder.Register(value)
	}
	iter.Release()
	err := iter.Error()
	if err == nil {
		db.prepared = true
	}
	return err
}

func (db *DB) setPrepCount(value int) error {
	key := []byte("prep-count")
	data := []byte(strconv.Itoa(value))
	return db.internal.Put(key, data, nil)
}

func (db *DB) incPrepCount(i int) error {
	size := db.prepCount()
	if size == -1 {
		return db.setPrepCount(1)
	}
	return db.setPrepCount(size + i)
}

func (db *DB) prepCount() int {
	err := db.Open()
	if err != nil {
		return -1
	}

	val, err := db.internal.Get([]byte("prep-count"), nil)
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(string(val))
	return n
}

// Checks if definition is present in current db.
func (db *DB) isPresent(def []byte) bool {
	err := db.Open()
	if err != nil {
		return false
	}

	key := []byte("prep:")
	key = append(key, def...)
	_, err = db.internal.Get(key, nil)
	return err == nil
}

// If not done already, registers type and example object in both
// the encoder and decoder.
func (db *DB) registerType(def, obj []byte) error {
	// Ensure that database is open.
	err := db.Open()
	if err != nil {
		return err
	}

	// Stop if type is already registered in database.
	if db.isPresent(def) {
		return nil
	}

	// Map prep#<n> to type definition bytes.
	k1 := []byte("prep#" + strconv.Itoa(db.prepCount()))
	err = db.incPrepCount(1)
	if err != nil {
		return err
	}
	v1 := []byte{}
	v1 = append(v1, def...)
	v1 = append(v1, obj...)
	err = db.internal.Put(k1, v1, nil)
	if err != nil {
		return err
	}

	// Map prep:<def> to empty string (for checking duplicate keys).
	k2 := []byte("prep:")
	k2 = append(k2, def...)
	v2 := []byte("")
	err = db.internal.Put(k2, v2, nil)
	if err != nil {
		db.internal.Delete(k1, nil)
		db.incPrepCount(-1)
		return err
	}
	return nil
}
