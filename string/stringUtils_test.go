package stringUtils

import (
	"math"
	"testing"
)


func TestReadableSize(t *testing.T ){
	s := HumanReadableSize( math.Pow(1024,4) + 900*math.Pow(1024,3))
	if s != "1.88 TB" {
		t.Errorf("Wrong number returned %s", s)
	}

	s = HumanReadableSize( math.Pow(1024,3) + 900*math.Pow(1024,2))
	if s != "1.88 GB" {
		t.Errorf("Wrong number returned %s", s)
	}

	s = HumanReadableSize( math.Pow(1024,2) + 900*math.Pow(1024,1))
	if s != "1.88 MB" {
		t.Errorf("Wrong number returned %s", s)
	}

	s = HumanReadableSize( math.Pow(1024,1) + 900)
	if s != "1.88 KB" {
		t.Errorf("Wrong number returned %s", s)
	}

	s = HumanReadableSize( float64(900 ))
	if s != "900 B" {
		t.Errorf("Wrong number returned %s", s)
	}

}
