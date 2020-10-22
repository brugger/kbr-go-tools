package fileUtils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}
}

func ReadIfFile(name string ) string {
	if IsFile( name ) {
		s, _ := ReadAll( name)
		return s
	}

	return name

}


func ReadAll(filename string ) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	return string( dat ), err
}

func Write( filename, data string, ) error {
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	return err
}

func IsFile(filename string ) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func FileSize(filename string ) int64 {
	if ! IsFile( filename ) {
		return -1
	} else {
	    s, err := os.Stat( filename )
	    checkErr( err )
	    return s.Size()
	}
}