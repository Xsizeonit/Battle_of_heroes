/*
 *  1. Creating XMLHttpRequest object that allow us to send ajax request
 *  2. Transfer data into json form that more comfortable
 *  3. Save json form of user data to use this data in future
 *  4. 244 - code if login and password is correct
 *     245 - code if user not exsist or inputing password not correct (error)
 *     246 - code if user try to input exsisting login in registation form (error)
 *     247 - code if user successfuly had registered
 *  5. If server has returned code 244 (successful log in) - open page main_table.html
 */
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
