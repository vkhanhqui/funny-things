var userId = null;
var websocket = null;
var receiver = null;

function init() {
	if ("WebSocket" in window) {
		userId = document.getElementsByName("userId")[0].value;
		while (userId === null) {
			userId = prompt("Enter userId");
		}
		document.getElementsByClassName("login-form")[0].style.display = "none";
		websocket = new WebSocket('ws://' + window.location.host + '/chat/' + userId);

		websocket.onopen = function(data) {
			document.getElementById("container").style.display = "block";
		};

		websocket.onmessage = function(data) {
			setMessage(JSON.parse(data.data));
		};

		websocket.onerror = function(e) {
			console.log('An error occured, closing application');

			cleanUp();
		};

		websocket.onclose = function(data) {
			cleanUp();

			var reason = (data.reason && data.reason !== null) ? data.reason : 'Goodbye';
			console.log(reason);
		};
	} else {
		console.log("Websockets not supported");
	}
}

function cleanUp() {
	document.getElementById("container").style.display = "none";
	userId = null;
	websocket = null;
}

function setReceiver(element) {
	receiver = element.id;
	console.log("receiver:" + receiver);
	document.getElementById("receiver").innerHTML = '<img src="http://' + window.location.host + '/static/images/chat_avatar_01.jpg"'
		+ 'alt="">'
		+ '<div>'
		+ '<br>'
		+ '<h2 id="receiver">' + receiver + '</h2>'
		+ '</div>';
	document.getElementById("chat").innerHTML = '';
}

function sendMessage() {
	var messageContent = document.getElementById("message").value;
	var message = buildMessage(userId, messageContent);

	document.getElementById("message").value = '';

	setMessage(message);
	console.log(message);
	websocket.send(JSON.stringify(message));
}

function buildMessage(userId, message) {
	return {
		userId: userId,
		message: message,
		receiver: receiver
	};
}

function setMessage(msg) {
	console.log("online users: " + msg.onlineList);
	if (msg.message != '[P]open' && msg.message != '[P]close') {
		var currentHTML = document.getElementById('chat').innerHTML;
		var newElem;

		if (msg.userId === userId) {
			newElem = '<li class="me">'
				+ '<div class="entete">'
				+ '<h3>10:12AM, Today</h3>'
				+ '<h2>' + msg.userId + '</h2>'
				+ '<span class="status blue"></span>'
				+ '</div>'
				+ '<div class="triangle"></div>'
				+ '<div class="message">' + msg.message + '</div>'
				+ '</li>';


		} else {
			newElem = '<li class="you">'
				+ '<div class="entete">'
				+ '<span class="status green"></span>'
				+ '<h2>' + msg.userId + '</h2>'
				+ '<h3>10:12AM, Today</h3>'
				+ '</div>'
				+ '<div class="triangle"></div>'
				+ '<div class="message">' + msg.message + '</div>'
				+ '</li>';

		}

		document.getElementById('chat').innerHTML = currentHTML
			+ newElem;
	} else {
		if (msg.message === '[P]open') {
			msg.onlineList.forEach(userId => setOnline(userId, true));
		} else {
			setOnline(msg.userId, false);
		}

	}
}

function setOnline(userId, is) {
	if (is === true) {
		document.getElementById('status-' + userId).innerHTML = '<span class="status green"></span>'
			+ 'online';
	} else {
		document.getElementById('status-' + userId).innerHTML = '<span class="status orange"></span>'
			+ 'offline';
	}
}