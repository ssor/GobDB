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
	"path"
	// "strconv"
	"fmt"
)

var g_debug = false

func debugOutput(args ...interface{}) {
	if g_debug {
		fmt.Println(args)
	}
}

// db path should like ./gobdb/people  , people is a folder which folds a list of people record
var default_db_root_path = "./gobdb"

func NewDB(name string, objPrototypeGenerator func() interface{}) *DB {
	if len(name) <= 0 {
		name = "demo"
	}

	location := path.Join(default_db_root_path, name)
	if objPrototypeGenerator == nil {
		// set a default string generator
		objPrototypeGenerator = func() interface{} {
			var s string
			return &s
		}
	}
	db := &DB{
		location:           location,
		Name:               name,
		ObjectsMap:         make(map[string]interface{}),
		prototypeGenerator: objPrototypeGenerator,
	}
	return db
}

type DB struct {
	location           string
	Name               string
	ObjectsMap         map[string]interface{}
	prototypeGenerator func() interface{}
}

//mkdir for db, read data if db already exits
func (db *DB) Init() (*DB, error) {
	// debugOutput("db location => ", db.location)
	if dry.FileExists(db.location) == false || dry.FileIsDir(db.location) == false {
		// debugOutput("location not exists: ", db.location)
		if err := os.MkdirAll(db.location, os.ModePerm); err != nil {
			return nil, err
		}
	}

	//read folder's files
	if files, err := dry.ListDirFiles(db.location); err != nil {
		return nil, err
	} else {
		// debugOutput("files: ", files)
		dry.StringEachMust(func(file string) {
			filePath := path.Join(db.location, file)
			obj := db.prototypeGenerator()
			if err := readFile(filePath, obj); err != nil {
				return
			}
			// debugOutput("gobdb => load file ", file, " ", obj)
			db.ObjectsMap[file] = obj

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

//if exists, delete, and create new file
func (db *DB) Update(key string, value interface{}) error {
	err := db.Delete(key)
	if err != nil {
		return err
	}
	return db.Put(key, value)
}

// Put encodes given key and value through gob
func (db *DB) Put(key string, value interface{}) error {
	filePath := path.Join(db.location, key)
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
	db.ObjectsMap[key] = value
	return nil
}

// Get ncodes given key via gob, fetches the corresponding
// value from within leveldb, and decodes that value into
// parameter two.
func (db *DB) Get(key string) interface{} {
	if _, ok := db.ObjectsMap[key]; ok {
		return db.ObjectsMap[key]
	} else {
		return nil
	}
}

// Has encodes given key via gob and checks if the resulting
// byte slice exists in the database's internal leveldb.
func (db *DB) Has(key string) bool {
	_, ok := db.ObjectsMap[key]
	return ok
}

// Delete encodes given key via gob, deleting the resulting
// byte slice from the database's internal leveldb.
func (db *DB) Delete(key string) error {
	if db.Has(key) == true {
		dbFilePath := path.Join(db.location, key)
		if dry.FileExists(dbFilePath) == false {
			return nil
		}
		if err := os.Remove(dbFilePath); err != nil {
			return err
		}
		delete(db.ObjectsMap, key)
	}
	return nil
}

// Entries counts key-value pairs in the database. This
// includes only pairs written through GobDB.Put.
func (db *DB) Count() int {
	return len(db.ObjectsMap)
}

func (db *DB) DB_FileExists(name string) bool {
	filePath := path.Join(db.location, name)
	return dry.FileExists(filePath)
}
