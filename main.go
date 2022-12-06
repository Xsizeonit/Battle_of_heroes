package main

import (
	"fmt"
	"encoding/json"
    "log"
    "net/http"
    "io/ioutil"
)

type test_struct struct {
    Test string
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//w.WriteHeader(http.StatusOK)
		fmt.Println("frontend" + r.URL.Path)
		if(r.URL.Path == "/") {
			http.ServeFile(w, r, "frontend/index.html")
		} else {
			http.ServeFile(w, r, "frontend" + r.URL.Path)
		}
	case "POST":
		w.WriteHeader(http.StatusOK)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		var t test_struct
		err = json.Unmarshal(body, &t)
		if err != nil {
			panic(err)
		}
		w.Write(body);
	default:
		fmt.Println("default")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	http.ListenAndServe(":3333", mux)
}
