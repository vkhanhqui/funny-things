
var username = null;
var websocket = null;
var receiver = null;

window.onload = function() {
	if ("WebSocket" in window) {
		username = document.getElementById("username").textContent;
		username = username.trim();
		websocket = new WebSocket('ws://' + window.location.host + '/chat/' + username);

		websocket.onopen = function() {
		};

		websocket.onmessage = function(data) {
			setMessage(JSON.parse(data.data));
		};

		websocket.onerror = function() {
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
	document.getElementsByClassName("container")[0].style.display = "none";
	username = null;
	websocket = null;
	receiver = null;
}

function setReceiver(element) {
	receiver = element.id;
	console.log("receiver: " + receiver);

	document.getElementById("receiver").innerHTML = '<div class="user-contact">' + '<div class="back">'
		+ '<i class="fa fa-arrow-left"></i>'
		+ '</div>'
		+ '<div class="user-contain">'
		+ '<div class="user-img">'
		+ '<img src="http://' + window.location.host + '/static/images/user-male.jpg"'
		+ 'alt="Image of user">'
		+ '<div class="user-img-dot"></div>'
		+ '</div>'
		+ '<div class="user-info">'
		+ '<span class="user-name">' + receiver + '</span>'
		+ '</div>'
		+ '</div>'
		+ '<div class="setting">'
		+ '<i class="fa fa-cog"></i>'
		+ '</div>' + '</div>'
		+ '<div class="list-messages-contain">'
		+ '<ul id="chat" class="list-messages">'
		+ '</ul>'
		+ '</div>'
		+ '<form class="form-send-message">'
		+ '<input id="message" type="text" class="txt-input" placeholder="Type message...">'
		+ '<label class="btn btn-image" for="attach"><i class="fa fa-file"></i></label>'
		+ ' <input type="file" id="attach"> <label class="btn btn-image" for="image"><i'
		+ ' class="fa fa-file-image-o"></i></label> <input type="file" id="image">'
		+ '<button type="button" class="btn btn-send" onclick="sendMessage();">'
		+ '<i class="fa fa-paper-plane"></i>'
		+ '</button>'
		+ '</form>';

	loadMessages(receiver);

	handleResponsive();
}

function handleResponsive() {
	var back = document.querySelector(".back");
	var rightSide = document.querySelector(".right-side");
	var leftSide = document.querySelector(".left-side");
	var conversation = document.querySelectorAll(".user-contain");

	if (back) {
		back.addEventListener("click", function() {
			rightSide.classList.remove("active");
			leftSide.classList.add("active");
		});
	}

	conversation.forEach(function(element, index) {
		element.addEventListener("click", function() {
			rightSide.classList.add("active");
			leftSide.classList.remove("active");
		});
	});
}

function sendMessage() {
	var rawData = document.getElementById("attach").files[0];
	var messageContent = document.getElementById("message").value;
	var messageType = "text";
	if (rawData == null) {
		document.getElementById("message").value = '';
	} else {
		document.getElementById("attach").value = null;
		console.log(rawData);
		messageContent = rawData.name;
		messageType = rawData.type;
	}
	var message = buildMessageToJson(username, messageContent, messageType);
	setMessage(message);
	console.log(message);
	websocket.send(JSON.stringify(message));
	if (rawData != null) {
		websocket.send(rawData);
	}
}

function buildMessageToJson(username, message, type) {
	return {
		username: username,
		message: message,
		type: type,
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
		goLastestMsg();
	} else {
		if (msg.message === '[P]open') {
			msg.onlineList.forEach(username => setOnline(username, true));
		} else {
			setOnline(msg.username, false);
		}

	}
}

function setOnline(username, isOnline) {
	var ele = document.getElementById('status-' + username);
	if (isOnline === true) {
		ele.classList.add('online');
	} else {
		ele.classList.remove('online');
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
			goLastestMsg();
		}
	};
	xhttp.open("GET", "http://" + window.location.host + "/chat-rest-controller?userId=" + userId, true);
	xhttp.send();
}

function customLoadMessage(sender, message) {
	var msgDisplay = '<li>'
		+ '<div class="message';
	if (username != sender) {
		msgDisplay += '">';
	} else {
		msgDisplay += ' right">';
	}
	return msgDisplay + '<div class="message-img">'
		+ '<img src="http://' + window.location.host + '/static/images/user-male.jpg" alt="">'
		+ ' </div>'
		+ '<div class="message-text">' + message + '</div>'
		+ '</div>'
		+ '</li>';
}

function goLastestMsg() {
	var msgLiTags = document.querySelectorAll(".message");
	var last = msgLiTags[msgLiTags.length - 1];
	last.scrollIntoView();
}