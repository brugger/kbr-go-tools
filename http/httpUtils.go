package httpUtils

import (
	"net/http"
)

func ArrayToMap(list []string) map[string]bool {
	m := make(map[string]bool)

	for _,l := range list {
		m[l] = true
	}

	return m
}



func AccessToken( r *http.Request ) string {

	//token := r.Header.Get("X-Session-Token")

	token := r.Header.Get("Authorization")
	if token == "" {
		return token
	}

	token = token[7:]

	return token

}



func ValidArguments( args map[string]interface{}, valid []string) bool {
	vm := ArrayToMap( valid)

	for k := range args {
		if _,ok := vm[ k ]; !ok {
			return false
		}
	}
	return true
}



func RequiredArguments( args map[string]interface{}, required []string) bool {

	for _,k := range required {
		if _,ok := args[ k ]; !ok {
			return false
		}
	}

	return true
}