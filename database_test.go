package GobDB

import (
	"testing"
    "os"
    "fmt"
    "io/ioutil"
    // "path/filepath"
  )


func TestBasic(t *testing.T) {
	// Make a database file.
	path, err := ioutil.TempDir("", "temp")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(path)


	// Make a database using that file.
	db := At(path)
	err = db.Open()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}


	// Map a key to a value.
	key := "name"
	val := "adam"
	err = db.Put(key, val)
	if err != nil {
		t.FailNow()
	}


	// Write, close, and reopen database (to check for persistance).
	db.Close()
	db.Open()


	// Fetch value from key.
	var out string
	err = db.Get(key, &out)
	if err != nil {
		t.FailNow()
	}


	// Check!
	if out != val {
		t.FailNow()
	}
}


func TestPersistence(t *testing.T) {
		// Make a database file.
	path, err := ioutil.TempDir("", "temp")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(path)


	// Make a database using that file.
	db1 := At(path)
	err = db1.Open()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}


	// Map a key to a value.
	key := "name"
	val := "adam"
	err = db1.Put(key, val)
	if err != nil {
		t.FailNow()
	}


	// Write, close, and reopen database (to check for persistance).
	db1.Close()
	db1.Open()


	// Fetch value from key.
	var out string
	err = db1.Get(key, &out)
	if err != nil {
		t.FailNow()
	}

	// Check!
	if out != val {
		t.FailNow()
	}

	db1.Close()

	// Make new db with same data file.
	db2 := At(path)
	db2.Open()

	out = ""
	err = db2.Get(key, &out)
	if err != nil {
		t.FailNow()
	}

	// Check!
	if out != val {
		t.FailNow()
	}

	db2.Close()
}