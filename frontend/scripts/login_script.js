function handleFormSubmit(event) {
	//Denie to reload page after button has pressed
	event.preventDefault();

	/*
	 * Save special tag for warning text into errorMes to in future display
	 * errors in user inputing. Plus reset warning text to delete previous
	 * error message.
	 */
	let errorMessage = document.getElementById('login_error_message');
	errorMessage.innerHTML = '';

	//Save input user data (login and password) into object named userInput
	let userInput = {
		Type: "login",
		Login: appForm.elements[0].value,
		Password: appForm.elements[1].value
	};
	
	console.log("Ok");

	/*
	 * Enter user data inro function prepareDataToSend(), that ckeck
	 * the presence of empty lines and spaces in the lines. It return
	 * object userInput with user data if all Ok (login and password is not empty and
	 * do not have spaces) and return false otherwise.
	 */
	userInput = prepareDataToSend(userInput);

	if(userInput != false) {
		sendData(userInput);
	}
	else {
		errorMessage.style.color = '#bd1217';
		errorMessage.innerHTML = 'Логин и пароль не могут быть пустыми и не могут содержать пробелы!';
	}
}

//Save login form into variable appForm
const appForm = document.getElementById('login_form');

//If button is pressed - go to the function handleFormSubmit
appForm.addEventListener('submit', handleFormSubmit);
