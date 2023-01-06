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

type user struct {
	User user_info
	InGame bool
	IsLogin bool
	Conn *websocket.Conn
	FriendConn *websocket.Conn
}

var list_users_to_websocket []user_info
var save_socket_users []*user

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
			list_users_to_websocket = append(list_users_to_websocket, user_input)
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
	log.Println("Create new socket")
	if (save_socket_users == nil) {
		save_socket_users = make([]*user, 0)
	}
	
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		r.Body.Close()
	}()
	
	conn, _ := upgrader.Upgrade(w, r, nil)
	
	ptr_users := &user{
		Conn: conn,
		IsLogin: false,
		InGame: false,
	}
	
	save_socket_users = append(save_socket_users, ptr_users)
	ptr_users.startThread()
}

func (i *user) startThread(){
	go func() {
		defer func() {
			var ind int;
			for index, ex_user := range save_socket_users {
				if(ex_user == i){
					ind = index;
					continue;
				}
				ex_user.Conn.WriteMessage(websocket.TextMessage, []byte("-" + i.User.Login))
			}
			i.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
			save_socket_users = append(save_socket_users[:ind], save_socket_users[ind+1:]...)
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		
		for {
			i.listen()
		}
	}()
}

func (i *user) listen() {
	_, b, err := i.Conn.ReadMessage()
	
	if err != nil {
		panic(err)
	}
	
	if(i.IsLogin == false) {
		var new_user user_info
		json.Unmarshal(b, &new_user)
		
		for index, ex_user := range list_users_to_websocket {
			if(new_user == ex_user) {
				i.User = new_user
				i.IsLogin = true
				list_users_to_websocket = append(list_users_to_websocket[:index], list_users_to_websocket[index+1:]...)
				log.Printf("New user with login: %s", i.User.Login)
				break
			}
		}
		for _, ex_user := range save_socket_users {
			if(ex_user == i) {
				continue
			}
			ex_user.Conn.WriteMessage(websocket.TextMessage, []byte(i.User.Login))
			i.Conn.WriteMessage(websocket.TextMessage, []byte(ex_user.User.Login))
		}
	} else {
		var friend_user *user
		right_user := false
		friend_login := string(b)
		for _, ex_user := range save_socket_users {
			if(friend_login == ex_user.User.Login) {
				friend_user = ex_user
				right_user = true
				break
			}
		}
		if(right_user == true) {
			log.Printf("User %s want to play with user %s", i.User.Login, friend_user.User.Login)
		}
	}
}

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
