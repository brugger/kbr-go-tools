package dbUtils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

var db_name string = "/tmp/go_db_test.db"

func TestMain(m *testing.M) {

	_, err := os.Stat(db_name)
	if !os.IsNotExist(err) {
		os.Remove(db_name)
	}

	Connect("sqlite3", db_name)
	err = Do("Create table TEST ( id  INTEGER PRIMARY KEY AUTOINCREMENT, idx int, value varchar(22));")
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error %s\n", err)
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestEscapeStringInt(t *testing.T) {
	s, _ := escapeString(int(1))

	if s != "1" {
		t.Errorf("Returned wrong value, expected \"1\" got %s ", s)
	}

}

func TestEscapeStringInt8(t *testing.T) {
	s, _ := escapeString(int8(1))

	if s != "1" {
		t.Errorf("Returned wrong value, expected \"1\" got %s ", s)
	}

}

func TestEscapeStringInt16(t *testing.T) {
	s, _ := escapeString(int16(1))

	if s != "1" {
		t.Errorf("Returned wrong value, expected \"1\" got %s ", s)
	}

}

func TestEscapeStringInt32(t *testing.T) {
	s, _ := escapeString(int32(1))

	if s != "1" {
		t.Errorf("Returned wrong value, expected \"1\" got %s ", s)
	}

}

func TestEscapeStringInt64(t *testing.T) {
	s, _ := escapeString(int64(1))

	if s != "1" {
		t.Errorf("Returned wrong value, expected \"1\" got %s ", s)
	}

}

func TestEscapeStringfloat32(t *testing.T) {
	s, _ := escapeString(float32(1.01))

	if s != "1.010000" {
		t.Errorf("Returned wrong value, expected \"1.01\" got %s ", s)
	}

}

func TestEscapeStringfloat64(t *testing.T) {
	s, _ := escapeString(float64(1.02))

	if s != "1.020000" {
		t.Errorf("Returned wrong value, expected \"1.01\" got %s ", s)
	}

}

func TestEscapeStringBoolTrue(t *testing.T) {
	s, _ := escapeString(true)

	if s != "true" {
		t.Errorf("Returned wrong value, expected \"true\" got %s ", s)
	}
}

func TestEscapeStringBoolFalse(t *testing.T) {
	s, _ := escapeString(false)

	if s != "false" {
		t.Errorf("Returned wrong value, expected \"false\" got %s ", s)
	}
}

func TestEscapeStringArray(t *testing.T) {
	ints := []int{1, 2, 3, 4}
	_, err := escapeString(ints)

	if err == nil {
		t.Errorf("accepted wrong input type %s", reflect.TypeOf(ints).String())
	}
}

func TestDo(t *testing.T) {
	//	t.Log("Insert data with Do")
	err := Do("INSERT INTO TEST (idx, value) values (1, 'hello'), (3, 'world');")
	checkErr(err)

}

func TestAdd(t *testing.T) {
	//	t.Log("Insert data with Add")
	var v = make(map[string]interface{})
	v["idx"] = 2
	v["value"] = "cruel"

	err := Add("TEST", v)
	checkErr(err)
}

func TestAddBulk(t *testing.T) {
	//	t.Log("Insert data with AddBulk")

	v1 := map[string]interface{}{"idx": 10, "value": "Remember to"}
	v2 := map[string]interface{}{"idx": 11, "value": "breathe"}

	err := AddBulk("TEST", []map[string]interface{}{v1, v2})
	checkErr(err)
}

func TestGetAll(t *testing.T) {
	v := GetAll("TEST")
	if len(v) != 5 {
		t.Errorf("Returned wrong number of entries expected 5 got %d", len(v))
	}
}

func TestGetById(t *testing.T) {
	v := GetByID("TEST", 2)

	if v["idx"].(int64) != 3 {
		t.Errorf("Returned wrong idx 3 got %d", v["idx"])
	}

}

func TestGetId(t *testing.T) {
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

func TestClose(t *testing.T) {
	Close()
}
