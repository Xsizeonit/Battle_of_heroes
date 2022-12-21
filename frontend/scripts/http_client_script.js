function sendData(data, errorFunc) {
	let xhr = new XMLHttpRequest();
	let jsonData = JSON.stringify(data);

	xhr.open("POST", "../../");
	xhr.setRequestHeader('Content-type', 'application/json; charset=utf-8');
	
	xhr.onreadystatechange = function() {
		if (xhr.readyState !== 4) return;
		
		if (xhr.status == 245)
			errorFunc("Неверно введён логин или пароль");
		if (xhr.status == 246)
			errorFunc("Этот логин уже используется")
		if (xhr.status == 247)
			errorFunc("Вы успешно зарегестрировались!");
	};
	xhr.send(jsonData);
}
