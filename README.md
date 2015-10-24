# Overview
GobDB is a simple database optimized for convenience in Go. It uses folder as table or collection,and gob file as row or atom.It's preformance may be lower than other db,like sqlite or bolt,but it's really convenience

GobDB 是一个为了方便在 Go 开发中的数据存储而开发的微型数据库，它的最大特点就是方便：
> 1. 使用文件作为存储方式，所以不依赖别的数据库
> 2. 存储的文件采用二进制序列化，所以占用空间很小
> 3. 仅有几个最简单的 API，一分钟就能上手
> 4. 最重要的是，代码简单，修改方便，有什么需求，稍微一动手就可以解决

最适合的场景：
> * 快速开发中不需要对数据库考虑太多
> * 数据量不多，对性能要求不高

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

Delete key
```
db.Delete("name")
```

You needn't close it,as it will save your change when you do that.


# Installation
```
go get github.com/ssor/GobDB
go test github.com/ssor/GobDB
```

#welcome pull request
欢迎一起开发，pull request me

