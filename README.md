# Overview
GobDB is a simple database optimized for convenience in Go. It wraps [leveldb](https://github.com/syndtr/goleveldb) to provide persistant key-value storage of [gob](http://golang.org/pkg/encoding/gob/)-compatible types.



# Sample Usage
Setup a database and assign a local data file.
```
db := GobDB.At("example")
db.Open()
```

Insert persistently key-value pairs. We use strings here, but all gob-compatible values are supported.
```
db.Put("name", "adam")
```

Fetch values of key-value pairs. Note that you must provide a pointer of the correct type. 
```
var value string = ""
db.Get("name", &value)
```

Check if keys are contained within the database.
```
db.Has("name")
db.Has("3234") 
```

Close the database and write changes to disk.
```
db.Close()
```



# Installation
```
go get github.com/dasmithii/GobDB
go test github.com/dasmithii/GobDB
```



# Miscellaneous Links
+ [GobDB GoDoc](http://godoc.org/github.com/dasmithii/GobDB)