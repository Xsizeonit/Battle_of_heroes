function equPasswordsInput(userPassword, userPasswordReplay) {
	if(userPassword === userPasswordReplay)
		return true;
	return false;
}

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
	if(obj.Password != appRegForm[2].value) {
		errorFunc("Пароли не совпадают");
		return false;
	}
		
	for(let key in obj){
		if(stringIsEmpty(obj[key]) == true) {
			errorFunc("Все поля должны быть заполненны");
			return false;
		}
		if(isSpaceInString(obj[key]) == true) {
			errorFunc("Логин и пароль не должны содержать пробелы");
			return false;
		}
	}
	return true;
}

function errorFunc(textError) {
	errorRegLabel.style.color = "#bd1217";
	errorRegLabel.innerHTML = textError;
}

function handleFormSubmit(event) {
	//Denie to reload page after button has pressed
	event.preventDefault();

	/*
	 * Save special tag for warning text into errorMes to in future display
	 * errors in user inputing. Plus reset warning text to delete previous
	 * error message.
	 */
	let errorMessage = document.getElementById('register_error_message');
	errorMessage.innerHTML = '';

	//Save input user data (login, password and password replay) into object named userInput
	let userInput = {
		Type: "registration",
		Login: appRegForm.elements[0].value,
		Password: appRegForm.elements[1].value,
	};


	/*
	 * Enter user data inro function prepareDataToSend(), that ckeck
	 * the presence of empty lines and spaces in the lines. It return
	 * object userInput with user data if all Ok (login and password is not empty and
	 * do not have spaces) and return false otherwise.
	 */
	let result = prepareDataToSend(userInput);

	if(result == true)
		sendData(userInput, errorFunc);
}

//Save login form into variable appForm
const appRegForm = document.getElementById("register_form");
let errorRegLabel = document.getElementById("register_error_message");

//If button is pressed - go to the function handleFormSubmit
appRegForm.addEventListener("submit", handleFormSubmit);
