const socket = new WebSocket("ws://192.168.0.103:3333/socket");
let userData;
let us = []
let btns = {}
socket.onopen = function(e) {
	userData = localStorage.getItem("data")
	console.log("userData: " + userData)
	socket.send(userData);
}

let container = document.getElementById("container");

socket.onmessage = function(e) {
	console.log("Get from server message");
	if((e.data[0] == "-")){
		let data = e.data.slice(1);
		container.removeChild(btns[data]);
		delete btns[data];
	}
	else if(us.includes(e.data) == false) {
		us.push(e.data)
		let new_btn = document.createElement("Button");
		new_btn.innerHTML = e.data;
		container.appendChild(new_btn);
		btns[e.data] = new_btn;
	}
}
