package httpUtils

import (
	"reflect"
	"testing"
)


func TestArrayToMap( t *testing.T) {

	a := []string{"A", "B", "C"}
	m := ArrayToMap( a )
	exp := map[string]bool{"A": true, "B": true, "C": true}
	if reflect.DeepEqual(m, exp) == false {
		t.Error(" wrong map returned!")
	}
}



func TestValidArguments( t *testing.T) {

	v := []string{"A", "B"}
	if ValidArguments( map[string]interface{}{"A":2, "B":"3"}, v) != true {
		t.Error("Arguments should be valid")
	}

	if ValidArguments( map[string]interface{}{"A":2, "C":"3"}, v) != false {
		t.Error("Arguments should be not valid")
	}
}



func TestRequiredArguments( t *testing.T) {

	v := []string{"A", "B"}

	if RequiredArguments( map[string]interface{}{"A":2, "B":"3"}, v) != true {
		t.Error("Arguments should be valid")
	}

	if RequiredArguments( map[string]interface{}{"A":2, "B":"3", "C":true}, v) != true {
		t.Error("Arguments should be valid")
	}

	if RequiredArguments( map[string]interface{}{"A":2, "C":true}, v) != false {
		t.Error("Arguments should be valid")
	}

}