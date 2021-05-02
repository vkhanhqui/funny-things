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
	document.getElementById("receiver").innerHTML = '<img src="http://' + window.location.host + '/static/images/chat_avatar_01.jpg"'
		+ 'alt="">'
		+ '<div>'
		+ '<br>'
		+ '<h2 id="receiver">' + receiver + '</h2>'
		+ '</div>';
	loadMessages(receiver);
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
		var currentChat = document.getElementById('chat').innerHTML;
		var newChatMsg = customLoadMessage(msg.username, msg.message);
		document.getElementById('chat').innerHTML = currentChat
			+ newChatMsg;
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

function loadMessages(userId) {
	var currentChatbox = document.getElementById("chat");
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			var messages = JSON.parse(this.responseText);
			var chatbox = "";
			messages.forEach(msg => {
				chatbox += customLoadMessage(msg.username, msg.message);
			});
			currentChatbox.innerHTML = chatbox;

			var senderLiTags = document.querySelectorAll(".me");
			var receiverLiTags = document.querySelectorAll(".you");
			var last = receiverLiTags[receiverLiTags.length - 1];
			if (senderLiTags.length == Math.max(senderLiTags.length, receiverLiTags)) {
				last = senderLiTags[senderLiTags.length - 1];
			}
			last.scrollIntoView();
		}
	};
	xhttp.open("GET", "http://" + window.location.host + "/chat-rest-controller?userId=" + userId, true);
	xhttp.send();
}

function customLoadMessage(sender, message) {
	if (username != sender) {
		return '<li class="you">'
			+ '<div class="entete">'
			+ '<span class="status green"></span>'
			+ '<h2>' + sender + '</h2>'
			+ '<h3>10:12AM, Today</h3>'
			+ '</div>'
			+ '<div class="triangle"></div>'
			+ '<div class="message">' + message + '</div>'
			+ '</li>';
	} else {
		return '<li class="me">'
			+ '<div class="entete">'
			+ '<h3>10:12AM, Today</h3>'
			+ '<h2>' + username + '</h2>'
			+ '<span class="status blue"></span>'
			+ '</div>'
			+ '<div class="triangle"></div>'
			+ '<div class="message">' + message + '</div>'
			+ '</li>';
	}
}