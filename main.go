package main

import (
	"encoding/json"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/websocket"
)

//Main information about user that connection to server
type user_info struct {
    Type string
    Login string
    Password string
}

/*
 *   All successful log in users server will be keeping in list_users_to_websocket
 *   to protectect server from unauthorized connecting to websocket. When
 *   user successful log in and create succussefuly websocket connection he
 *   send to server again his data (login and password) with using websocket connection
 *   to confirm his data.
 */
var list_users_to_websocket []user_info

func parse_post_request(w http.ResponseWriter, r *http.Request) {
	//Trying to connect to mysql database
	db, err := sql.Open("mysql", "root:root@/USERS")
	if err != nil {
		panic(err)
	}
	
	//Close connect to database when function pase_post_request() have finished
	defer db.Close()
	
	//Give and transform camming data from json form into struct user_info
	var user_input user_info
	json.NewDecoder(r.Body).Decode(&user_input)
	
	switch user_input.Type {
	case "login":
		/*
		 *   Get password that corresponds inputing user login.
		 *   If user has inputed right password - send back code 244 (successful log in)
		 *   and add user to list users for websocket connection.
		 *   Else send back code 245 (not right inputing password).
		 */
		row, _ := db.Query("select hash_password from users where login=\"" + user_input.Login + "\"")
		row.Next()
		
		var user_right_password string
		row.Scan(&user_right_password)
		
		if(user_right_password == user_input.Password) {
			w.WriteHeader(244)
			list_users_to_websocket = append(list_users_to_websocket, user_input)
		} else {
			w.WriteHeader(245)
		}
	case "registration":
		/*
		 *   Check absense inputing registration login in database.
		 *   If inputing login has existed in database - send code 246 (users
		 *   with this login is existing).
		 *   Else add new user in database and send back 247 code (successful registation)
		 */
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


/*
 *   Analyse user requests
 *   If he send post requets then we riese parse_post_request function.
 * 
 *   If user want to connect to server websocket then server rise spical function
 *   parse_socket to do this. All websocket methods is placed into websocket.go file.
 * 
 *   Else if user connect to main page (to /) then we give him index.html
 *   file that keep login form.
 * 
 *   Else if user want to get other files (other html files or scripts and styles)
 *   server send to client this files.
 */
func home(w http.ResponseWriter, r *http.Request) {
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
