package main

import (
	"encoding/json"
    "net/http"
    "log"
    "github.com/gorilla/websocket"
)

/*
 *   Main info about users for websocket connect:
 *   field User keeps login and password about user
 * 
 *   field InGame indicates that user in game and not ready to accept
 *   requests from other users
 * 
 *   field InLogin idicates avaiability server info about user login and 
 *   password in websocket connection
 * 
 *   filed Conn keep websocket connection with current user
 * 
 *   field FriendConn keep connection to other user (need for game and 
 *   empty if user is not in game)
 */
type user struct {
	User user_info
	InGame bool
	IsLogin bool
	WantPlay bool
	Conn *websocket.Conn
	FriendUser *user
}

var save_socket_users []*user

//Allow cross domain requests
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

//This function will be risening when new user want to connect to websocket
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
	
	//Create connection (transofrm http request into websocket request)
	conn, _ := upgrader.Upgrade(w, r, nil)
	
	//Create new user
	ptr_users := &user{
		Conn: conn,
		IsLogin: false,
		InGame: false,
		WantPlay: false,
	}
	
	//Save new user and start new thread
	save_socket_users = append(save_socket_users, ptr_users)
	ptr_users.startThread()
}

func (i *user) startThread(){
	//Start function in new thread with gorutiens
	go func() {
		/*
		 *   When user is disconecting then server send all user information that
		 *   other user clients delete disconecting user from game menu
		 */
		defer func() {
			var ind int;
			/*
			 *   Send all users info about disconect user and parallel
			 *   find index this user in array of users
			 */
			for index, ex_user := range save_socket_users {
				if(ex_user == i){
					ind = index;
					continue;
				}
				ex_user.Conn.WriteMessage(websocket.TextMessage, []byte("-" + i.User.Login))
			}
			
			//Close connection anyway
			i.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
			
			log.Printf("Delete user %s", save_socket_users[ind].User.Login)
			
			//delete user from array
			save_socket_users = append(save_socket_users[:ind], save_socket_users[ind+1:]...)
			
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()
		
		//Listen message from user
		for {
			i.listen()
		}
	}()
}

func (i *user) listen() {
	//Read message
	_, b, err := i.Conn.ReadMessage()
	
	if err != nil {
		panic(err)
	}
	
	/*
	 *   If user has only just websocket connected (IsLogin = false) then
	 *   he must then his log in data again that server may identify this user
	 *   from user websocket user connection
	 */
	if(i.IsLogin == false) {
		//Get data about user
		var new_user user_info
		json.Unmarshal(b, &new_user)
		
		//Find user to list of log in users and delete this user from it
		for index, ex_user := range list_users_to_websocket {
			if(new_user == ex_user) {
				i.User = new_user
				i.IsLogin = true
				
				//Delete founded user from list
				list_users_to_websocket = append(list_users_to_websocket[:index], list_users_to_websocket[index+1:]...)
				log.Printf("New user with login: %s", i.User.Login)
				break
			}
		}
		
		/*
		 *   Send all users (except new user) info about new user and
		 *   send current new user info about other users
		 */
		for _, ex_user := range save_socket_users {
			if(ex_user == i) {
				continue
			}
			ex_user.Conn.WriteMessage(websocket.TextMessage, []byte(i.User.Login))
			i.Conn.WriteMessage(websocket.TextMessage, []byte(ex_user.User.Login))
		}
	} else if(i.WantPlay == false){
		/*
		 *   If user already has loginned that it may only send other user
		 *   request to play with him
		 */
		i.WantPlay = true
		//var friend_user *user
		right_user := false
		//friend_user_accept := false
		
		//Get nickname other user
		friend_login := string(b)
		
		//Find this user from exsiting users
		for _, ex_user := range save_socket_users {
			if(friend_login == ex_user.User.Login) {
				//friend_user = ex_user
				i.FriendUser = ex_user
				ex_user.FriendUser = i
				right_user = true
				break
			}
		}
		
		/*
		 *   If other user has founded then we assing (save in special field - 
		 *   FriendConn connection other user)
		 */
		if(right_user == true) {
			i.FriendUser.WantPlay = true
			i.FriendUser.Conn.WriteMessage(websocket.TextMessage, []byte("+" + i.User.Login))
			
			//i.FriendConn = friend_user.Conn
			//i.FriendUser = i
			log.Printf("Now user %s and user %s trying to connect", i.User.Login, i.FriendUser.User.Login)
		}
	} else if(i.WantPlay == true) {
		ans := string(b)
		if(ans == "1") {
			log.Printf("User %s ready to fight with user %s", i.User.Login, i.FriendUser.User.Login)
			i.FriendUser.Conn.WriteMessage(websocket.TextMessage, []byte("@" + i.FriendUser.User.Login))
			i.WantPlay = false
			i.FriendUser.WantPlay = false
		} else if(ans == "0"){
			log.Printf("User %s NOT ready to fight with user %s", i.User.Login, i.FriendUser.User.Login)
			i.FriendUser.Conn.WriteMessage(websocket.TextMessage, []byte("#" + i.FriendUser.User.Login))
			i.WantPlay = false
			i.FriendUser.WantPlay = false
		}
	}
}
