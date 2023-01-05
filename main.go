package main

import (
	"encoding/json"
    "net/http"
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/websocket"
)

type user_info struct {
    Type string
    Login string
    Password string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func parse_post_request(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@/USERS")
	if err != nil {
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
			w.WriteHeader(244)
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

func parse_socket(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade: ", err);
		return
	}

	defer conn.Close()
	
	for {
		messageType, message, err := conn.ReadMessage()
		
		if err != nil {
			log.Print("Error during connection reading: ", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Print("Error during connection writing: ", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.Path)
	switch r.Method {
	case "GET":
		if(r.URL.Path == "/socket") {
			parse_socket(w, r)
		} else if(r.URL.Path == "/") {
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
