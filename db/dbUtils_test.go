package dbUtils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"errors"
)
var db_name string = "/tmp/go_db_test.db"


func TestMain(m *testing.M) {

	_, err := os.Stat(db_name)
	if !os.IsNotExist(err) {
		os.Remove( db_name)
	}

	Connect("sqlite3", db_name)
	err = Do("Create table TEST ( id  INTEGER PRIMARY KEY AUTOINCREMENT, idx int, value varchar(22));")
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf( "Error %s\n", err )
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestDo ( t *testing.T ) {
//	t.Log("Insert data with Do")
	err := Do("INSERT INTO TEST (idx, value) values (1, 'hello'), (3, 'world');")
	checkErr( err )

}

func TestAdd ( t *testing.T ) {
//	t.Log("Insert data with Add")
	var v = make(map[string]interface{})
	v["idx"] = 2
	v["value"] = "cruel"

	err := Add("TEST", v)
	checkErr( err )
}

func TestAddBulk ( t *testing.T ) {
//	t.Log("Insert data with AddBulk")

	v1 := map[string]interface{}{"idx": 10, "value": "Remember to"}
	v2 := map[string]interface{}{"idx": 11, "value": "breathe"}

	err := AddBulk("TEST", []map[string]interface{}{v1,v2})
	checkErr( err )
}

func TestGetAll( t *testing.T ) {
	v := GetAll("TEST")
	if len(v) != 5 {
		t.Errorf("Returned wrong number of entries expected 5 got %d", len( v ))
	}
}

func TestGetById( t *testing.T ) {
	v := GetByID("TEST", 2)

	if v["idx"].(int64) != 3 {
		t.Errorf("Returned wrong idx 3 got %d", v["idx"])
	}

}

func TestGetId( t *testing.T ) {
	id := GetID("TEST", map[string]interface{}{"value": "world"})
	//fmt.Println( id )

	if id.(int64) != 2 {
		t.Errorf("Returned wrong id, expected 2 got %d", id)
	}
}


func TestUpdate(t *testing.T) {
	v := GetByID("TEST", 2)
	v["value"] = "New value"
	Update("TEST", v, map[string]interface{}{"id": 2})
	v = GetByID("TEST", 2)
	//fmt.Println( v )
}


func TestDelete(t *testing.T) {
	Delete("TEST", 2)
	v := GetByID("TEST", 2)
	if len(v) != 0 {
		t.Errorf("Element not deleted")
	}
	//fmt.Println( v )
}


func TestClose( t *testing.T) {
	Close()
}
