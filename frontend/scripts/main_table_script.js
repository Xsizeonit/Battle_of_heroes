function main_table() {
	//Creating and saving socket into localStorage that using it in future
	const socket = new WebSocket("ws://192.168.0.103:3333/socket");
	localStorage.setItem("socket", socket);
	
	const jsonUserData = localStorage.getItem("data");
	
	//Create constant cantainer to add buttons into menu in future
	const container = document.getElementById("container");
	
	/* 
	 *  Array users_nicknames is keeping users nicknames :)
	 *  Objeckt buttons is keeping buttons, that allow current user connecting to other users
	 */
	let users_nicknames = [];
	let buttons = {};
	
	/*
	 *  Create and save in localStorage variable InGame that indicates
	 *  denieded access other users (except friend/rival) to current user
	 */
	let inGame = false;
	localStorage.setItem("inGame", inGame);
	
	//If socket has opened - sending users data to server that it has all info about user
	socket.onopen = (openEvent) => {
		socket.send(jsonUserData);
	}
	
	socket.onmessage = (onmessageEvent) => {
		if(onmessageEvent.data == "100") {
			inGame = true;
			//Game
		} else if((inGame == false) && (onmessageEvent.data[0] == "-")) {
			
			/*
			 *   If user is not playing and the first character of message is "-" (without quotes)
			 *   that some users has disconected and script must delete this user from menu
			 */
			 
			 //Get user login without character "-"
			let data = onmessageEvent.data.slice(1);
			
			//Delete button with user nickname from menu
			container.removeChild(buttons[data]);
			
			//Delete button from array buttons
			delete buttons[data];
			
			//Delete nickname disconected user from array nicknames
			users_nicknames.forEach(function(nickname, index, users_nicknames) {
				if(nickname == data) {
					users_nicknames.splice(index, 1);
					return;
				}
			});
		} else if((inGame == false) && (users_nicknames.includes(onmessageEvent.data) == false)) {
			
			/*
			 *   If user not playing and new user do not visiable at this time
			 *   then script add button with new user nickname
			 */
			 
			//Add new user nickname to array nicknames
			users_nicknames.push(onmessageEvent.data)
			
			//Create new button and add user nickname as text on this button
			const new_btn = document.createElement("Button");
			new_btn.innerHTML = onmessageEvent.data;
			
			//Add (and visiable) button to menu
			container.appendChild(new_btn);
			
			buttons[onmessageEvent.data] = new_btn;
			
			/*
			 *   If user has clicked on button then socket send then
			 *   socket will send button text (selected user login) to
			 *   server
			 */
			buttons[onmessageEvent.data].onclick = function() {
				socket.send(onmessageEvent.data)
			}
		}
	}
}

main_table();
