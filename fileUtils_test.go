package dbUtils

import (
	"testing"
)


func TestIsFile(t *testing.T ){
	if ! IsFile( "/etc/passwd") {
		t.Error("file is there /etc/passwd")
	}
	if  IsFile( "/etc/sdfsdfsdfsdfsf"){
		t.Error("file should not exist", )
	}

}
