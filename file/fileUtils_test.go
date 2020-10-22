package fileUtils

import (
	"testing"
)

var tmp_file string = "/tmp/go_text_file.txt"
var file_content string = "There was a boy\nA very strange enchanted boy"

func TestIsFile(t *testing.T ){
	if ! IsFile( "/etc/passwd") {
		t.Error("file is there /etc/passwd")
	}
	if  IsFile( "/etc/sdfsdfsdfsdfsf"){
		t.Error("file should not exist", )
	}

}


func TestFileSize(t *testing.T) {
	size := FileSize("/etc/passwd")
	if size < 0 {
		t.Error("No size found for /etc/passwd")
	}
}


func TestWrite(t *testing.T) {
	Write(tmp_file, file_content)
}

func TestReadAll(t *testing.T) {
	s, _ := ReadAll(tmp_file)
	if s != file_content {
		t.Errorf("file content did not match expectation %s", s)
	}
}

func TestReadIfFile(t *testing.T) {
	s := ReadIfFile(tmp_file)
	if s != file_content {
		t.Errorf("file content did not match expectation %s", s)
	}

	s = ReadIfFile("Not-a-File")
	if s != "Not-a-File" {
		t.Errorf("file content did not match expectation %s", s)
	}
}
