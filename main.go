package main

import (
	"fmt"
	"encoding/json"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type user_info struct {
    Type string
    Login string
    Password string
}

func parse_post_request(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@/USERS")
	if err != nil {
		fmt.Println("Error")
		panic(err)
	}
	defer db.Close()

	var user_input user_info
	json.NewDecoder(r.Body).Decode(&user_input)
	
	switch user_input.Type {
	case "login":
		row, _ := db.Query("select hash_password from users where login=\"" + user_input.Login + "\"")
		row.Next()
		
		var user_right_password string
		row.Scan(&user_right_password)
		
		if(user_right_password == user_input.Password) {
			fmt.Println("Ok")
			//Websocket
		} else {
			w.WriteHeader(245)
		}
	case "registration":
		row, _ := db.Query("select login from users where login=\"" + user_input.Login + "\"")
		row.Next()
		
		var user_exsist_login string
		row.Scan(&user_exsist_login)
		
		if(user_input.Login == user_exsist_login) {
			w.WriteHeader(246)
		} else {
			w.WriteHeader(247)
			db.Exec("insert into users (login, hash_password, score) values (?, ?, ?)",
			user_input.Login, user_input.Password, 0)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if(r.URL.Path == "/") {
			http.ServeFile(w, r, "frontend/index.html")
		} else {
			http.ServeFile(w, r, "frontend" + r.URL.Path)
		}
	case "POST":
		parse_post_request(w, r)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	http.ListenAndServe(":3333", mux)
}
