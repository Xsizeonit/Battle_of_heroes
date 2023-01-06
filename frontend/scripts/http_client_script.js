
function sendData(data, errorFunc) {
	let xhr = new XMLHttpRequest();
	let jsonData = JSON.stringify(data);
	localStorage.setItem("data", jsonData)
	
	xhr.open("POST", "../../");
	xhr.setRequestHeader('Content-type', 'application/json; charset=utf-8');
	
	xhr.onreadystatechange = function() {
		if (xhr.readyState !== 4) return;
		switch (xhr.status) {
			case 244:
				window.location.href = "../main_table.html";
				break;
			case 245:
				errorFunc("Неверно введён логин или пароль");
				break;
			case 246:
				errorFunc("Этот логин уже используется");
				break;
			case 247:
				errorFunc("Вы успешно зарегестрировались!");
				break;
		}
	};
	xhr.send(jsonData);
}
