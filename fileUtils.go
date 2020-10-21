package dbUtils


import (
	"os"
)


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