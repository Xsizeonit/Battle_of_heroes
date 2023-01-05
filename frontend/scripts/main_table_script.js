const appLoginForm = document.getElementById("form");
let socket = new WebSocket("ws://192.168.0.103:3333/socket");

appLoginForm.addEventListener("submit", handleFormSubmit);

function handleFormSubmit(e) {
	event.preventDefault();
	let inp = document.getElementById("input")

	if(socket.readyState == 1) {
		socket.send(inp.value);
	}

}
