package endpoints

import (
	"encoding/json"
	"net/http"
)

func PublicEntry(w http.ResponseWriter, r *http.Request) {
	converted, err := json.Marshal("Welcome to my project! check README or the PROG2005 git wiki for legal URLS")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(converted)
}
