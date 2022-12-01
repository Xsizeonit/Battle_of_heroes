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


function sendData(data, serverURL) {
	let xhr = new XMLHttpRequest();
	let json_data = JSON.stringify(data);
	
	xhr.open("POST", serverURL);
	
	xhr.setRequestHeader('Content-type', 'application/json; charset=utf-8');
	
	xhr.send(json_data);
	
	if (xhr.status != 200) {
		console.log(xhr.status + ': ' + xhr.statusText);
	} else {
		console.log(xhr.responseText);
	}
}
