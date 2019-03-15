package handler

import (
	"net/http"
)

func Beat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//data,_:=json.Marshal(utils.ValueFromContext(r.Context(),utils.RUserkey))
	w.Write([]byte("alive"))
}
