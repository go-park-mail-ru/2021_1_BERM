package apiserver

import (
	"fmt"
	"net/http"
)



func setupSimpleResponse(w *http.ResponseWriter, req *http.Request) {
	print("1: ")
	fmt.Println(req.Header.Values("Origin"))
	print("2: ")
	fmt.Println(req.Header.Values("origin"))
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	if cookie := req.Cookies(); len(cookie) != 0 {
		(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func setupDifficultResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Max-Age", "6400")
}
