# Overview
GobDB is a simple database optimized for convenience in Go. It wraps [leveldb](https://github.com/basho/leveldb) to provide persistant key-value storage of [gob](http://golang.org/pkg/encoding/gob/)-compatible types.



# Example
```
package main
import "github.com/dasmithii/GobDB"

func main() {
	// Set up a database and assign it a local data file.
	db := GobDB.At("example")
	db.Open()
	defer db.Close()

	// Insert persistantly-stored key-value pairs. We use strings
	// here, but all gob-compatible values are supported.
	db.Put("name", "adam")

	// Fetch values of key-value pairs. Note that you must provide
	// a pointer of the correct type. 
	var value string
	db.Get(name, &value)

	Check if keys are contained within the database.
	db.Has("name") // => true
	db.Has("3234") // => false
}
```



# Installation
```
go get github.com/dasmithii/GobDB
```