package controller

import "net/http"

func Subscribe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
