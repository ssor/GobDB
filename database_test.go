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



// func TestStructs(t *testing.T) {}