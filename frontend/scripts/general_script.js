function stringIsEmpty(str) {
	if(str === '')
		return true;
	return false;
}

function isSpaceInString(str) {
	for(let index in str)
		if(str[index] === " ")
			return true;
	return false;
}

//If string is empry or contains at least one space - error, return false in calling function
function prepareDataToSend(obj) {
	for(let key in obj)
		if(stringIsEmpty(obj[key]) || isSpaceInString(obj[key]))
			return false;
	return obj;
}


function sendData(data) {
	let xhr = new XMLHttpRequest();
	let json_data = JSON.stringify(data);

	xhr.open("POST", "../../");
	xhr.setRequestHeader('Content-type', 'application/json; charset=utf-8');
	
	xhr.onreadystatechange = function() {
		if (xhr.readyState !== 4) return;
		if (xhr.status == 245) {
			document.getElementById("login_error_message").style.color = "red"
			document.getElementById("login_error_message").innerHTML = "Неверно введён логин или пароль"
		}
		if (xhr.status == 246) {
			document.getElementById("register_error_message").style.color = "red"
			document.getElementById("register_error_message").innerHTML = "Этот логин уже используется"
		}
		if (xhr.status == 247) {
			document.getElementById("register_error_message").style.color = "red"
			document.getElementById("register_error_message").innerHTML = "Вы успешно зарегестрировались!"
		}
	};
	xhr.send(json_data);
}
