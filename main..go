package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Nieless/tpro/tpro"
	"github.com/gorilla/mux"
	"net/http"
)

var tineUrlData = make(map[string]string)

func main()  {

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/shorting", ShortingUrlHandler).Methods("POST")
	myRouter.HandleFunc("/shorting", GetLongUrl).Methods("GET")

	err := tpro.ListenAndServe("8888", myRouter)
	if err != nil {
		panic(err)
	}
}

// ShortingUrlHandler serves the html file
func ShortingUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	type shortingUrl struct {
		Url string `json:"url"`
	}

	su := &shortingUrl{}
	_, err := tpro.DecodeJSON(r, su)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytesUrl := []byte(su.Url)
	hashedUrlData := fmt.Sprintf("%x", md5.Sum(bytesUrl))
	tinyUrl := hashedUrlData[0:6]
	tineUrlData[tinyUrl] = su.Url

	if err := json.NewEncoder(w).Encode(tinyUrl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetLongUrl(w http.ResponseWriter, r *http.Request) {
	tinyUrl := r.FormValue("tiny_url")
	fmt.Println(tineUrlData)
	longUrl := tineUrlData[tinyUrl]
	fmt.Println(tinyUrl)

	http.Redirect(w, r, longUrl, http.StatusSeeOther)
}
