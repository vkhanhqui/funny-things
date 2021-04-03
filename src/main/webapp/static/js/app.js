var email = null;
var websocket = null;

function init() {
	if ("WebSocket" in window) {
		email = document.getElementsByName("email")[0].value;
		while (email === null) {
			email = prompt("Enter email");
		}
		document.getElementsByClassName("login-form")[0].style.display = "none";
		websocket = new WebSocket('ws://localhost:8080/chat-controller/' + email);
		websocket.onopen = function(data) {
			document.getElementById("main").style.display = "block";
		};

		websocket.onmessage = function(data) {
			setMessage(JSON.parse(data.data));
		};

		websocket.onerror = function(e) {
			alert('An error occured, closing application');

			cleanUp();
		};

		websocket.onclose = function(data) {
			cleanUp();

			var reason = (data.reason && data.reason !== null) ? data.reason : 'Goodbye';
			alert(reason);
		};
	} else {
		alert("Websockets not supported");
	}
}

function cleanUp() {
	document.getElementById("main").style.display = "none";

	email = null;
	websocket = null;
}

function sendMessage() {
	var messageContent = document.getElementById("message").value;
	var message = buildMessage(email, messageContent);

	document.getElementById("message").value = '';

	setMessage(message);
	websocket.send(JSON.stringify(message));
}

function buildMessage(email, message) {
	return {
		email: email,
		message: message
	};
}

function setMessage(msg) {
	var currentHTML = document.getElementById('scrolling-messages').innerHTML;
	var newElem;

	if (msg.email === email) {
		newElem = '<p style="background: #ebebe0;"><span>' + msg.email
			+ ' : ' + msg.message + '</span></p>';
	} else {
		newElem = '<p><span>' + msg.email + ' : ' + msg.message
			+ '</span></p>';
	}

	document.getElementById('scrolling-messages').innerHTML = currentHTML
		+ newElem;
}