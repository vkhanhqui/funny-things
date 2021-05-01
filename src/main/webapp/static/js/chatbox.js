var username = null;
var websocket = null;
var receiver = null;

function init() {
	if ("WebSocket" in window) {
		username = document.getElementById("username").textContent;
		websocket = new WebSocket('ws://' + window.location.host + '/chat/' + username);

		websocket.onopen = function(data) {

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
	username = null;
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
	var message = buildMessage(username, messageContent);

	document.getElementById("message").value = '';

	setMessage(message);
	console.log(message);
	websocket.send(JSON.stringify(message));
}

function buildMessage(username, message) {
	return {
		username: username,
		message: message,
		receiver: receiver
	};
}

function setMessage(msg) {
	console.log("online users: " + msg.onlineList);
	if (msg.message != '[P]open' && msg.message != '[P]close') {
		var currentHTML = document.getElementById('chat').innerHTML;
		var newElem;

		if (msg.username === username) {
			newElem = '<li class="me">'
				+ '<div class="entete">'
				+ '<h3>10:12AM, Today</h3>'
				+ '<h2>' + msg.username + '</h2>'
				+ '<span class="status blue"></span>'
				+ '</div>'
				+ '<div class="triangle"></div>'
				+ '<div class="message">' + msg.message + '</div>'
				+ '</li>';


		} else {
			newElem = '<li class="you">'
				+ '<div class="entete">'
				+ '<span class="status green"></span>'
				+ '<h2>' + msg.username + '</h2>'
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
			msg.onlineList.forEach(username => setOnline(username, true));
		} else {
			setOnline(msg.username, false);
		}

	}
}

function setOnline(username, is) {
	if (is === true) {
		document.getElementById('status-' + username).innerHTML = '<span class="status green"></span>'
			+ 'online';
	} else {
		document.getElementById('status-' + username).innerHTML = '<span class="status orange"></span>'
			+ 'offline';
	}
}