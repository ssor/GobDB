package GobDB

import (
	"fmt"
	// "io/ioutil"
	// "github.com/ungerik/go-dry"
	"os"
	"testing"
	// "path/filepath"
)

var dbRootPath = "./gobdb"

// var dbPath = dbRootPath + "/example"
var dbExample = "example"

func TestInit(t *testing.T) {
	//clear env
	clearEnv := func() {
		os.RemoveAll(dbRootPath)
	}
	clearEnv()
	defer clearEnv()

	_, err := NewDB(dbExample, nil).Init()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	db, err := NewDB(dbExample, nil).Init()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if db == nil {
		t.FailNow()
	}
}

type ExampleThing struct {
	Name string
	Age  int
}

func (e *ExampleThing) Say() string {
	return fmt.Sprintf("I am %s, and %d years old", e.Name, e.Age)
}
func NewExample(name string, age int) *ExampleThing {
	return &ExampleThing{
		Name: name,
		Age:  age,
	}
}

func TestStructs(t *testing.T) {
	clearEnv := func() {
		os.RemoveAll(dbRootPath)
	}
	clearEnv()
	defer clearEnv()

	//create db
	db, _ := NewDB(dbExample, nil).Init()
	err := db.Put("first", NewExample("first", 1))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if db.Has("first") == false {
		t.Log("first not exists")
		t.FailNow()
	}
	//test if file created
	if db.DB_FileExists("first") == false {
		t.Log("file not exists")
		t.FailNow()
	}
	err = db.Put("second", NewExample("second", 2))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// reload db
	db2, err := NewDB(dbExample, func() interface{} {
		var exmaple ExampleThing
		return &exmaple
	}).Init()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if db2.Count() != 2 {
		t.Log("count error")
		t.FailNow()
	}

	//get obj
	second := db2.Get("second")
	if second == nil {
		t.Log("load db failed")
		t.FailNow()
	}
	if _, ok := second.(*ExampleThing); !ok {
		t.Log("type assert error")
		t.FailNow()
	}

	//delete obj
	err = db2.Delete("first")
	if err != nil {
		t.Log("delete obj error")
		t.FailNow()
	}
	if db2.Has("first") == true {
		t.Log("first should be deleted")
		t.FailNow()
	}
	if db2.Count() != 1 {
		t.Log("records count should be 1")
		t.FailNow()
	}
	if db2.DB_FileExists("first") == true {
		t.Log("file should not exists")
		t.FailNow()
	}

}
