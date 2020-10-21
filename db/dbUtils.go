package dbUtils

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"reflect"
)

var db *sql.DB

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}
}

func Connect(driver string, db_uri string) {

	dbConn, err := sql.Open(driver, db_uri)
	checkErr(err)
	db = dbConn
	return
}

func Close() {
	err := db.Close()
	checkErr( err )
}

func escapeString(v interface{}) string {

	switch reflect.TypeOf( v ).String() {
	case "string":
		return fmt.Sprintf("'%s'", v)
	case "int", "int16", "int32", "int64":
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("'%s'", v)
	}
}

func asList(query string) []map[string]interface{} {

	rows, err := db.Query(query)
	checkErr(err)
	cols, err := rows.Columns() // Remember to check err afterwards
	checkErr(err)

//	fmt.Println(cols)

	var res []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			checkErr(err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		res = append(res, m)
	}

	return res
}

func Get(table string, filter map[string]interface{}, rest ...string) []map[string]interface{} {

	stmt := fmt.Sprintf("SELECT * FROM %s ", table)

	var conds []string
	for key, value := range filter {
		switch reflect.TypeOf( value ).String() {
			case "string":
				conds = append(conds, fmt.Sprintf(" %s = '%s'", key, value))
			case "int":
				conds = append(conds, fmt.Sprintf(" %s = %d", key, value))
			default:
				conds = append(conds, fmt.Sprintf(" %s = '%s'", key, value))
		}
	}


	//    fmt.Println( conds )
	if len(conds) > 0 {
		stmt = fmt.Sprintf("%s WHERE %s ", stmt, strings.Join(conds[:], " AND "))
	}
	if len(rest) > 0 {
		stmt = fmt.Sprintf("%s %s", stmt, strings.Join(rest[:], " "))
	}

	return asList(stmt)
}

func Do(query string ) error {
	_, err := db.Exec( query )
	return err
}


func GetSingle( table string, filter map[string]interface{}, rest ...string) map[string]interface{} {

//	fmt.Println(reflect.TypeOf(rest).String())

	values := Get(table, filter, rest...)
	if len( values ) > 1 {
		fmt.Println("Function returned multiple vales!")
		return nil
	} else if len( values ) == 1 {
		return values[0]
	} else {
		return nil
	}

}

func GetAll( table string, rest ...string) []map[string]interface{} {
	return Get(table, make(map[string]interface{}), rest...)
}

func GetByID( table string, id interface{}, rest ...string) map[string]interface{} {
	f := make(map[string]interface{})
	f["id"] = id
	return GetSingle(table, f, rest...)
}


func GetID( table string, filter map[string]interface{}, rest ...string) interface{} {

	//	fmt.Println(reflect.TypeOf(rest).String())

	value := GetSingle(table, filter, rest...)
	if value != nil{
		return value["id"]
	}

	return nil
}


func Add( table string, entry map[string]interface{}) error {


	var keys, values  []string

	for k := range entry {
		//fmt.Printf(" VALUE %s --> %s\n", k, entry[k])
		keys = append( keys, k)
		values = append( values, escapeString(entry[k]))

	}

	stmt := fmt.Sprintf( "INSERT INTO %s (%s) VALUES ( %s )", table, strings.Join( keys, ","), strings.Join( values, ","))


	err := Do( stmt )

	return err
}


func AddBulk( table string, entries []map[string]interface{}) error {

	var all_values[]string
	var all_keys string

	for _, entry := range entries {
		var keys, values  []string

		for k := range entry {
			keys = append( keys, k)
			values = append( values, escapeString(entry[k]))
		}
		all_keys = strings.Join( keys, ",")
		all_values = append( all_values, strings.Join( values, ","))


	}
	stmt := fmt.Sprintf( "INSERT INTO %s (%s) VALUES ( %s )", table, all_keys, strings.Join(all_values, "), ( "))
	//fmt.Println( stmt )

	err := Do( stmt )
	return err
}


func Update(table string, values map[string]interface{}, conditions map[string]interface{}) {


	var updates, conds  []string

	for k := range values {
		updates = append( updates, fmt.Sprintf("%s = %s", k, escapeString(values[k])))
	}

	for k := range conditions {
		conds = append( conds, fmt.Sprintf("%s = %s", k, escapeString(values[k])))
	}

	stmt := fmt.Sprintf( "UPDATE %s SET %s WHERE %s ", table, strings.Join( updates, ", "), strings.Join( conds, " AND "))

	//fmt.Println( stmt )
	err := Do( stmt )
	checkErr( err )

}

func Delete( table string, id interface{}) {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = %s", table, escapeString(id))
	err := Do(stmt)
	checkErr( err )

}
