// Package GobDB implements a persistant key-value store of
// gob-compatible types. This is accomplished with a light
// wrapper around leveldb and Go's gob encoding library.
package GobDB

import (
	// "github.com/syndtr/goleveldb/leveldb"
	// "github.com/syndtr/goleveldb/leveldb/util"
	"encoding/gob"
	"errors"
	"github.com/ungerik/go-dry"
	"os"
	// "strconv"
	// "fmt"
)

func NewDB(path string, objPrototypeGenerator func() interface{}) *DB {
	if objPrototypeGenerator == nil {
		// set a default string generator
		objPrototypeGenerator = func() interface{} {
			var s string
			return &s
		}
	}
	db := &DB{
		location:           path,
		Name:               path,
		objectIndex:        make(map[string]interface{}),
		prototypeGenerator: objPrototypeGenerator,
	}
	return db
}

// db path should like ./gobdb/people  , people is a folder which folds a list of people record

// var default_db_path = "./gobdb/"

// DB is a LevelDB wrapper that stores key-value pairs of
// gob-compatible types.
type DB struct {
	location           string
	Name               string
	objectIndex        map[string]interface{}
	prototypeGenerator func() interface{}
}

//mkdir for db, read data if db already exits
func (db *DB) Init() (*DB, error) {
	// if dry.FileExists(default_db_path) == false || dry.FileIsDir(default_db_path) == false { //may be it's a file
	// 	if err := os.Mkdir(default_db_path, os.O_CREATE); err != nil {
	// 		return err
	// 	}
	// }
	if dry.FileExists(db.location) == false || dry.FileIsDir(db.location) == false {
		if err := os.MkdirAll(db.location, os.ModePerm); err != nil {
			return nil, err
		}
	}

	//read folder's files
	if files, err := dry.ListDirFiles(db.location); err != nil {
		return nil, err
	} else {
		dry.StringEach(func(file string) error {
			filePath := db.location + "/" + file
			obj := db.prototypeGenerator()
			if err := readFile(filePath, obj); err != nil {
				return err
			}
			// fmt.Println("load file ", file, " ", obj)
			db.objectIndex[file] = obj
			return nil
		}, files)
	}
	return db, nil
}

// read and decode file
func readFile(path string, obj interface{}) error {
	// fmt.Println(obj)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err2 := dec.Decode(obj)
	// fmt.Println(obj)

	if err2 != nil {
		return err2
	}
	return nil
}

// Put encodes given key and value through gob
func (db *DB) Put(key string, value interface{}) error {
	filePath := db.location + "/" + key
	if dry.FileExists(filePath) == true {
		return errors.New("file already exists")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(value)
	if err != nil {
		return err
	}
	db.objectIndex[key] = value
	return nil
}

// Get ncodes given key via gob, fetches the corresponding
// value from within leveldb, and decodes that value into
// parameter two.
func (db *DB) Get(key string) interface{} {
	if _, ok := db.objectIndex[key]; ok {
		return db.objectIndex[key]
	} else {
		return nil
	}
}

// Has encodes given key via gob and checks if the resulting
// byte slice exists in the database's internal leveldb.
func (db DB) Has(key string) bool {
	_, ok := db.objectIndex[key]
	return ok
}

// Delete encodes given key via gob, deleting the resulting
// byte slice from the database's internal leveldb.
func (db *DB) Delete(key string) error {
	if db.Has(key) == true {
		if err := os.Remove(db.location + "/" + key); err != nil {
			return err
		}
		delete(db.objectIndex, key)
	}
	return nil
}

// Entries counts key-value pairs in the database. This
// includes only pairs written through GobDB.Put.
func (db *DB) Count() int {
	return len(db.objectIndex)
}
