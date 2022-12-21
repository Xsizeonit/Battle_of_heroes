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
	errorLoginLabel.style.color = "#bd1217";
	errorLoginLabel.innerHTML = textError;
}

function handleFormSubmit(event) {
	//Denie to reload page after button has pressed
	event.preventDefault();

	/*
	 * Save special tag for warning text into errorMes to in future display
	 * errors in user inputing. Plus reset warning text to delete previous
	 * error message.
	 */
	errorLoginLabel.innerHTML = "";

	//Save input user data (login and password) into object named userInput
	let userInput = {
		Type: "login",
		Login: appLoginForm[0].value,
		Password: appLoginForm[1].value
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
const appLoginForm = document.getElementById("login_form");
let errorLoginLabel = document.getElementById("login_error_message");

//If button is pressed - go to the function handleFormSubmit
appLoginForm.addEventListener("submit", handleFormSubmit);
