package dbUtils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
	checkErr(err)
}

func escapeString(v any) (string, error) {

	switch reflect.TypeOf(v).String() {
	case "string":
		return fmt.Sprintf("'%s'", v), nil
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("%d", v), nil
	case "float32", "float64":
		return fmt.Sprintf("%f", v), nil
	case "bool":
		return fmt.Sprintf("%t", v), nil
	default:
		return "", errors.New(fmt.Sprintf("Unsupported type passed %s", reflect.TypeOf(v).String()))
	}
}

func AsList(query string) ([]map[string]interface{}, error) {

	if query[len(query)-1] != ';' {
		query = fmt.Sprintf("%s;", query)
	}

	fmt.Println(query)
	rows, err := db.Query(query)
	checkErr(err)
	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		return nil, nil
	}

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

	return res, nil
}

func Get(table string, filter map[string]interface{}, rest ...string) ([]map[string]interface{}, error) {

	stmt := fmt.Sprintf("SELECT * FROM %s ", table)

	var conds []string
	for key, value := range filter {
		v, err := escapeString(value)
		if err != nil {
			var e []map[string]interface{}
			return e, err
		}
		conds = append(conds, fmt.Sprintf(" %s = %s", key, v))
	}

	//    fmt.Println( conds )
	if len(conds) > 0 {
		stmt = fmt.Sprintf("%s WHERE %s ", stmt, strings.Join(conds[:], " AND "))
	}
	if len(rest) > 0 {
		stmt = fmt.Sprintf("%s %s", stmt, strings.Join(rest[:], " "))
	}

	// as list will crash if input is illegal
	return AsList(stmt)
}

func Do(query string) error {
	if query[len(query)-1] != ';' {
		query = fmt.Sprintf("%s;", query)
	}
	fmt.Println(query)
	_, err := db.Exec(query)
	return err
}

func GetSingle(table string, filter map[string]interface{}, rest ...string) map[string]interface{} {

	//	fmt.Println(reflect.TypeOf(rest).String())

	values, err := Get(table, filter, rest...)
	if err != nil {
		return nil
	}
	if len(values) > 1 {
		fmt.Println("Function returned multiple vales!")
		return nil
	} else if len(values) == 1 {
		return values[0]
	} else {
		return nil
	}

}

func GetAll(table string, rest ...string) ([]map[string]interface{}, error) {
	return Get(table, make(map[string]interface{}), rest...)
}

func GetByID(table string, id interface{}, rest ...string) map[string]interface{} {
	f := make(map[string]interface{})
	f["id"] = id
	return GetSingle(table, f, rest...)
}

func GetID(table string, filter map[string]interface{}, rest ...string) interface{} {

	//	fmt.Println(reflect.TypeOf(rest).String())

	value := GetSingle(table, filter, rest...)
	if value != nil {
		return value["id"]
	}

	return nil
}

func Add(table string, entry map[string]interface{}) error {

	var keys, values []string

	for k := range entry {
		//fmt.Printf(" VALUE %s --> %s\n", k, entry[k])
		keys = append(keys, k)
		value, err := escapeString(entry[k])
		if err != nil {
			return err
		}
		values = append(values, value)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES ( %s )", table, strings.Join(keys, ","), strings.Join(values, ","))

	err := Do(stmt)

	return err
}

func AddBulk(table string, entries []map[string]interface{}) error {

	var all_values []string
	var all_keys string

	for _, entry := range entries {
		var keys, values []string

		for k := range entry {
			keys = append(keys, k)
			value, err := escapeString(entry[k])
			if err != nil {
				return err
			}
			values = append(values, value)
		}
		all_keys = strings.Join(keys, ",")
		all_values = append(all_values, strings.Join(values, ","))

	}
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES ( %s )", table, all_keys, strings.Join(all_values, "), ( "))
	//fmt.Println( stmt )

	err := Do(stmt)
	return err
}

func Update(table string, values map[string]interface{}, conditions map[string]interface{}) {

	var updates, conds []string

	for k := range values {
		value, err := escapeString(values[k])
		checkErr(err)
		updates = append(updates, fmt.Sprintf("%s = %s", k, value))
	}

	for k := range conditions {
		value, err := escapeString(conditions[k])
		checkErr(err)
		conds = append(conds, fmt.Sprintf("%s = %s", k, value))
	}

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s ", table, strings.Join(updates, ", "), strings.Join(conds, " AND "))

	//fmt.Println( stmt )
	err := Do(stmt)
	checkErr(err)

}

func Delete(table string, id interface{}) {
	value, err := escapeString(id)
	checkErr(err)

	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = %s", table, value)
	err = Do(stmt)
	checkErr(err)

}
