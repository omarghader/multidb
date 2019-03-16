package multidb

import (
	"fmt"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/omarghader/multidb"
)

func TestArangoDatabase(t *testing.T) {
	dbName := "test1"

	db := NewArangodb(multidb.ConnectionOptions{
		Host:     "localhost",
		Port:     "859",
		Username: "root",
		Password: "root",
		DBName:   dbName,
	})

	if db.Create() == nil {
		t.Errorf("Creation should return an error")
	}

	db = NewArangodb(multidb.ConnectionOptions{
		Host:     "localhost",
		Port:     "8529",
		Username: "root",
		Password: "root",
		DBName:   dbName,
	})

	if db == nil {
		t.Errorf(multidb.EXCEPTION_CONNECTION_FAILED)
		return
	}

	if db.Exists() {
		err := db.Drop()
		if err != nil {
			t.Errorf("Cannot drop database")
			return
		}
	}

	err := db.Create()
	if err != nil {
		t.Errorf("Cannot create database")
		return
	}

	db.Drop()

}

func TestArangoTable(t *testing.T) {
	dbName := "test1"

	db := NewArangodb(multidb.ConnectionOptions{
		Host:     "localhost",
		Port:     "8529",
		Username: "root",
		Password: "root",
		DBName:   dbName,
	})

	db.Create()

	tbl := db.Table("table_test")
	if tbl.Exists() {
		err := tbl.Drop()
		if err != nil {
			t.Errorf("Cannot drop table")
		}
	}

	err := tbl.Create()
	if err != nil {
		t.Errorf("Cannot create table")
	}

	doc := map[string]string{"_key": "doc1", "name": "document1"}
	_, err = tbl.Insert("doc1", doc)
	if err != nil {
		t.Errorf("Cannot insert table")
	}

	doc1, err := tbl.Find("doc1")
	if err != nil {
		t.Errorf("Cannot find table")
	}

	if doc1.(map[string]interface{})["name"] != "document1" {
		t.Errorf("Cannot find object from table")
	}

	doc["name"] = "document2"
	_, err = tbl.Update("doc1", doc)
	if err != nil {
		t.Errorf("Cannot update table")
	}

	doc2, err := tbl.Find("doc1")
	if err != nil {
		t.Errorf("Cannot find table")
	}

	if doc2.(map[string]interface{})["name"] != "document2" {
		t.Errorf("Cannot updae object from table")
	}

	_, err = tbl.Delete("doc1")
	if err != nil {
		t.Errorf("Cannot find table")
	}

	tbl.Drop()
	db.Drop()
}

type DocTest struct {
	Key   string `json:"_key"`
	Id    string `json:"_id"`
	Name  string
	Age   int
	Likes []string
}

func TestArangoGraph(t *testing.T) {
	dbName := "test1"
	db := NewArangodb(multidb.ConnectionOptions{
		Host:     "localhost",
		Port:     "8529",
		Username: "root",
		Password: "root",
		DBName:   dbName,
	})

	// db.Drop()
	db.Create()

	db.Table("test1").Create()

	var d1, d2 DocTest
	d1.Name = "Diego"
	d1.Age = 22
	d1.Likes = []string{"arangodb", "golang", "linux"}
	d1.Key = "diego"
	d1.Id = "test1/diego"

	d1Res, err := db.Table("test1").Insert("", d1)
	if err != nil {
		panic(err)
	}

	d2.Name = "Facundo"
	d2.Age = 25
	d2.Likes = []string{"php", "linux", "python"}
	d2.Key = "facundo"
	d2.Id = "test1/facundo"

	d2Res, err := db.Table("test1").Insert("", d2)
	if err != nil {
		panic(err)
	}

	var graphName = "graphTest"
	graph := db.Graph(graphName, []string{"test1"}, []string{"test1"})
	graph.Create()

	rel := graph.Relation("is_related")
	fmt.Printf("Connecting %s => %s \n", d1.Id, d2.Id)

	relation, err := rel.Insert((d1Res.(driver.DocumentMeta)).ID.String(),
		(d2Res.(driver.DocumentMeta)).ID.String(), map[string]interface{}{"is": "friend"})
	if err != nil {
		t.Errorf("Relation creation error : %s\n", err)
	}

	if relation == nil {
		t.Errorf("Relation should not be nil\n")
	}

	_, err = rel.Insert((d2Res.(driver.DocumentMeta)).ID.String(),
		(d2Res.(driver.DocumentMeta)).ID.String(), map[string]interface{}{"is": "me"})
	if err != nil {
		t.Errorf("Relation creation error : %s\n", err)
	}

	// err = db.Drop()
	if err != nil {
		panic(err)
	}
}
