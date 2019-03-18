package arangodb

import (
	"errors"
	"fmt"
	"testing"

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
			t.Errorf("Cannot drop table: %s", err)
			panic(err)
		}
	}

	err := tbl.Create()
	if err != nil {
		t.Errorf("Cannot create table")
	}

	doc := DocTest{
		Key:  "doc1",
		Name: "document1",
	}
	var docRes Document
	_, err = tbl.Insert(doc, &docRes)
	if err != nil {
		t.Errorf("Cannot insert table")
	}

	var doc1 DocTest
	_, err = tbl.Find(docRes.Key, &doc1)
	if err != nil {
		t.Errorf("Cannot find table %s", err)
		panic(err)
	}

	if doc1.Name != "document1" {
		t.Errorf("Cannot find object from table")
	}

	doc.Name = "document2"
	_, err = tbl.Update(docRes.Key, doc, nil)
	if err != nil {
		t.Errorf("Cannot update table %s", err)
	}

	var doc2 DocTest
	_, err = tbl.Find(docRes.Key, &doc2)
	if err != nil {
		t.Errorf("Cannot find table")
	}

	if doc2.Name != "document2" {
		t.Errorf("Cannot updae object from table")
	}

	_, err = tbl.Delete(docRes.Key, nil)
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

	var d1Res Document
	_, err := db.Table("test1").Insert(d1, &d1Res)
	if err != nil {
		panic(err)
	}

	if d1Res.Key != d1.Key {
		panic(errors.New("Insert : result does not match input "))
	}

	d2.Name = "Facundo"
	d2.Age = 25
	d2.Likes = []string{"php", "linux", "python"}
	d2.Key = "facundo"
	d2.Id = "test1/facundo"

	var d2Res Document
	_, err = db.Table("test1").Insert(d2, &d2Res)
	if err != nil {
		panic(err)
	}

	var graphName = "graphTest"
	graph := db.Graph(graphName)
	graph.Create()

	rel := graph.Relation("is_related")
	fmt.Printf("Connecting %s => %s \n", d1.Id, d2.Id)

	relation, err := rel.Insert(
		Edge{
			From: d1Res.ID,
			To:   d2Res.ID,
		},
		map[string]interface{}{"is": "friend"})
	if err != nil {
		t.Errorf("Relation creation error : %s\n", err)
	}

	if relation == nil {
		t.Errorf("Relation should not be nil\n")
	}

	_, err = rel.Insert(

		Edge{
			From: d2Res.ID,
			To:   d2Res.ID,
		}, map[string]interface{}{"is": "me"})
	if err != nil {
		t.Errorf("Relation creation error : %s\n", err)
	}

	// err = db.Drop()
	if err != nil {
		panic(err)
	}
}
