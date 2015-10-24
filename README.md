# Overview
GobDB is a simple database optimized for convenience in Go. It uses folder as table or collection,and file as row or atom.It's preformance may be lower than other db,like sqlite or bolt,but it's really convenience

# Sample Usage
Setup a database and assign a local data file.
```
db := GobDB.New("example").Init()

```

Insert persistently key-value pairs. We use strings here, but all gob-compatible values are supported.
```
db.Put("name", "adam")
```

Fetch values of key-value pairs. Note that you will get the correct type. 
```
value := db.Get("name")
```

Check if keys are contained within the database.
```
db.Has("name")
db.Has("3234") 
```

You needn't close it,as it will save your change when you do that.


# Installation
```
go get github.com/ssor/GobDB
go test github.com/ssor/GobDB
```


