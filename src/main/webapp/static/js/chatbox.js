
var username = null;
var websocket = null;
var receiver = null;
var userAvatar = null;
var receiverAvatar = null;

var back = null;
var rightSide = null;
var leftSide = null;
var conversation = null;

var attachFile = null;
var imageFile = null;
var file = null;
var listFile = [];
var typeFile = "image";
var deleteAttach = null;

window.onload = function() {
	if ("WebSocket" in window) {
		username = document.getElementById("username").textContent;
		userAvatar = document.getElementById("userAvatar").textContent;
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
			console.log(data.reason);
			cleanUp();
		};
	} else {
		console.log("Websockets not supported");
	}
}

function cleanUp() {
	username = null;
	websocket = null;
	receiver = null;
}

handleResponsive();

function setReceiver(element) {
	receiver = element.id;
	console.log("receiver: " + receiver);
	receiverAvatar = document.getElementById('img-' + receiver).src;
	var status = '';
	if (document.getElementById('status-' + receiver).classList.contains('online')) {
		status = 'online';
	}
	var rightSide = '<div class="user-contact">' + '<div class="back">'
		+ '<i class="fa fa-arrow-left"></i>'
		+ '</div>'
		+ '<div class="user-contain">'
		+ '<div class="user-img">'
		+ '<img src="' + receiverAvatar + '" '
		+ 'alt="Image of user">'
		+ '<div class="user-img-dot ' + status + '"></div>'
		+ '</div>'
		+ '<div class="user-info">'
		+ '<span class="user-name">' + receiver + '</span>'
		+ '</div>'
		+ '</div>'
		+ '<div class="setting">'
		+ '<i class="fa fa-cog"></i>'
		+ '</div>'
		+ '</div>'
		+ '<div class="list-messages-contain">'
		+ '<ul id="chat" class="list-messages">'
		+ '</ul>'
		+ '</div>'
		+ '<form class="form-send-message">'
		+ '<ul class="list-file"></ul> '
		+ '<input type="text" id="message" class="txt-input" placeholder="Type message...">'
		+ '<label class="btn btn-image" for="attach"><i class="fa fa-file"></i></label>'
		+ '<input type="file" multiple id="attach">'
		+ '<label class="btn btn-image" for="image"><i class="fa fa-file-image-o"></i></label>'
		+ '<input type="file" accept="image/*" multiple id="image">'
		+ '<button type="button" class="btn btn-send" onclick="sendMessage();">'
		+ '<i class="fa fa-paper-plane"></i>'
		+ '</button>'
		+ '</form>';

	document.getElementById("receiver").innerHTML = rightSide;
	
	loadMessages();

	handleResponsive();

	displayFiles();
	
	makeFriend(rightSide);
}

function makeFriend(rightSide) {
	fetch("http://" + window.location.host + "/friend-rest-controller?sender=" + username + "&receiver=" + receiver)
		.then(function(data) {
			return data.json();
		})
		.then(data => {
			console.log(data);
			if (data.status == false && data.owner == username) {
				rightSide = '<div class="user-contact">' + '<div class="back">'
					+ '<i class="fa fa-arrow-left"></i>'
					+ '</div>'
					+ '<div class="user-contain">'
					+ '<div class="user-img">'
					+ '<img src="' + receiverAvatar + '" '
					+ 'alt="Image of user">'
					+ '<div class="user-img-dot ' + status + '"></div>'
					+ '</div>'
					+ '<div class="user-info">'
					+ '<span class="user-name">' + receiver + '</span>'
					+ '</div>'
					+ '</div>'
					+ '<div class="setting">'
					+ '<i class="fa fa-cog"></i>'
					+ '</div>'
					+ '<form action="http://localhost:8080/chat" method="post" >'
					+ '<input type="hidden" name="sender" value="' + username + '">'
					+ '<input type="hidden" name="receiver" value="' + receiver + '">'
					+ '<input type="hidden" name="status" value="false">'
					+ '<input type="hidden" name="isAccept" value="false">'
					+ '<input type="submit" value="Cho KB">'
					+ '</form>'
					+ '</div>'
					+ '<div class="list-messages-contain">'
					+ '<ul id="chat" class="list-messages">'
					+ '</ul>'
					+ '</div>';
					
					document.getElementById("receiver").innerHTML = rightSide;
			} else if (data.status == false && data.owner != username) {
				rightSide = '<div class="user-contact">' + '<div class="back">'
					+ '<i class="fa fa-arrow-left"></i>'
					+ '</div>'
					+ '<div class="user-contain">'
					+ '<div class="user-img">'
					+ '<img src="' + receiverAvatar + '" '
					+ 'alt="Image of user">'
					+ '<div class="user-img-dot ' + status + '"></div>'
					+ '</div>'
					+ '<div class="user-info">'
					+ '<span class="user-name">' + receiver + '</span>'
					+ '</div>'
					+ '</div>'
					+ '<div class="setting">'
					+ '<i class="fa fa-cog"></i>'
					+ '</div>'
					+ '<form action="http://localhost:8080/chat" method="post" >'
					+ '<input type="hidden" name="sender" value="' + username + '">'
					+ '<input type="hidden" name="receiver" value="' + receiver + '">'
					+ '<input type="hidden" name="status" value="true">'
					+ '<input type="hidden" name="isAccept" value="true">'
					+ '<input type="submit" value="Accept New Friend">'
					+ '</form>'
					+ '</div>'
					+ '<div class="list-messages-contain">'
					+ '<ul id="chat" class="list-messages">'
					+ '</ul>'
					+ '</div>';
					document.getElementById("receiver").innerHTML = rightSide;
			}

		})
		.catch(ex => {
			rightSide = '<div class="user-contact">' + '<div class="back">'
				+ '<i class="fa fa-arrow-left"></i>'
				+ '</div>'
				+ '<div class="user-contain">'
				+ '<div class="user-img">'
				+ '<img src="' + receiverAvatar + '" '
				+ 'alt="Image of user">'
				+ '<div class="user-img-dot ' + status + '"></div>'
				+ '</div>'
				+ '<div class="user-info">'
				+ '<span class="user-name">' + receiver + '</span>'
				+ '</div>'
				+ '</div>'
				+ '<div class="setting">'
				+ '<i class="fa fa-cog"></i>'
				+ '</div>'
				+ '<form action="http://localhost:8080/chat" method="post" >'
				+ '<input type="hidden" name="sender" value="' + username + '">'
				+ '<input type="hidden" name="receiver" value="' + receiver + '">'
				+ '<input type="hidden" name="status" value="false">'
				+ '<input type="hidden" name="isAccept" value="false">'
				+ '<input type="submit" value="Add Friend">'
				+ '</form>'
				+ '</div>'
				+ '<div class="list-messages-contain">'
				+ '<ul id="chat" class="list-messages">'
				+ '</ul>'
				+ '</div>';

			document.getElementById("receiver").innerHTML = rightSide;
		});
}

function handleResponsive() {
	back = document.querySelector(".back");
	rightSide = document.querySelector(".right-side");
	leftSide = document.querySelector(".left-side");
	conversation = document.querySelectorAll(".user-contain");

	if (back) {
		back.addEventListener("click", function() {
			rightSide.classList.remove("active");
			leftSide.classList.add("active");
			listFile = [];
			renderFile();
		});
	}

	conversation.forEach(function(element, index) {
		element.addEventListener("click", function() {
			rightSide.classList.add("active");
			leftSide.classList.remove("active");
		});
	});
}

function displayFiles() {
	attachFile = document.getElementById("attach");
	imageFile = document.getElementById("image");
	file = document.querySelector(".list-file");
	deleteAttach = document.querySelectorAll(".delete-attach");

	attachFile.addEventListener("change", function(e) {
		let filesInput = e.target.files;

		for (const file of filesInput) {
			listFile.push(file);
		}

		typeFile = "file";
		renderFile("attach");

		this.value = null;
	});

	imageFile.addEventListener("change", function(e) {
		let filesImage = e.target.files;

		for (const file of filesImage) {
			listFile.push(file);
		}

		typeFile = "image";

		renderFile("image");

		this.value = null;
	});



}

function deleteFile(idx) {
	if (!isNaN(idx)) listFile.splice(idx, 1);

	renderFile(typeFile);
}

function renderFile(typeFile) {
	let listFileHTML = "";
	let idx = 0;

	if (typeFile == "image") {
		for (const file of listFile) {
			listFileHTML += '<li><img src="' + URL.createObjectURL(file)
				+ '" alt="Image file"><span data-idx="'
				+ (idx) + '" onclick="deleteFile('
				+ idx + ')" class="delete-attach">X</span></li>';
			idx++;
		}
	} else {
		for (const file of listFile) {
			listFileHTML += '<li><div class="file-input">' + file.name
				+ '</div><span data-idx="'
				+ (idx) + '" onclick="deleteFile('
				+ idx + ')" class="delete-attach">X</span></li>';
			idx++;
		}
	}


	if (listFile.length == 0) {
		file.innerHTML = "";
		file.classList.remove("active");
	} else {
		file.innerHTML = listFileHTML;
		file.classList.add("active");
	}

	deleteAttach = document.querySelectorAll(".delete-attach");
}

function sendMessage() {
	var inputText = document.getElementById("message").value;
	if (inputText != '') {
		sendText();
	} else {
		sendAttachments();
	}
}

function sendText() {
	var messageContent = document.getElementById("message").value;
	var messageType = "text";
	document.getElementById("message").value = '';
	var message = buildMessageToJson(messageContent, messageType);
	setMessage(message);
	console.log(message);
	websocket.send(JSON.stringify(message));
}

function sendAttachments() {
	var messageType = "attachment";
	for (file of listFile) {
		messageContent = file.name.trim();
		messageType = file.type;
		var message = buildMessageToJson(messageContent, messageType);
		console.log(message);
		websocket.send(JSON.stringify(message));
		websocket.send(file);
		message.message = '<img src="' + URL.createObjectURL(file) + '" alt="">';
		if (messageType.startsWith("audio")) {
			message.message = '<audio controls>'
				+ '<source src="' + URL.createObjectURL(file) + '" type="' + messageType + '">'
				+ '</audio>';
		} else if (messageType.startsWith("video")) {
			message.message = '<video width="400" controls>'
				+ '<source src="' + URL.createObjectURL(file) + '" type="' + messageType + '">'
				+ '</video>';
		}
		setMessage(message);
	}
	file = document.querySelector(".list-file");
	file.classList.remove("active");
	file.innerHTML = "";
	listFile = [];
}

function buildMessageToJson(message, type) {
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
			msg.onlineList.forEach(username => {
				try {
					setOnline(username, true);
				} catch (ex) { }
			});
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

function loadMessages() {
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
	xhttp.open("GET", "http://" + window.location.host + "/chat-rest-controller?sender=" + username
		+ "&receiver=" + receiver, true);
	xhttp.send();
}

function customLoadMessage(sender, message) {
	var imgSrc = receiverAvatar;
	var msgDisplay = '<li>'
		+ '<div class="message';
	if (username != sender) {
		msgDisplay += '">';
	} else {
		imgSrc = userAvatar;
		msgDisplay += ' right">';
	}
	return msgDisplay + '<div class="message-img">'
		+ '<img src="' + imgSrc + '" alt="">'
		+ ' </div>'
		+ '<div class="message-text">' + message + '</div>'
		+ '</div>'
		+ '</li>';
}

function searchFriendByKeyword(keyword) {
	fetch("http://" + window.location.host + "/users-rest-controller?username=" + username + "&keyword=" + keyword)
		.then(function(data) {
			return data.json();
		})
		.then(data => {
			console.log(data);
			document.querySelector(".list-user").innerHTML = "";
			data.forEach(function(data) {
				if (data.isOnline) status = "online";
				else status = "";

				let appendUser = '<li id="' + data.username + '" onclick="setReceiver(this);">'
					+ '<div class="user-contain">'
					+ '<div class="user-img">'
					+ '<img id="img-' + data.username + '"'
					+ ' src="http://' + window.location.host + '/files/' + data.username + '/' + data.avatar + '"'
					+ 'alt="Image of user">'
					+ '<div id="status-' + data.username + '" class="user-img-dot ' + status + '"></div>'
					+ '</div>'
					+ '<div class="user-info">'
					+ '<span class="user-name">' + data.username + '</span>'
					+ '<span'
					+ '</div>'
					+ '</div>'
					+ '</li>';
				document.querySelector(".list-user").innerHTML += appendUser;
			});
		});

	//
}

function searchUser(ele) {
	searchFriendByKeyword(ele.value);
	console.log(ele.target);
}

function goLastestMsg() {
	var msgLiTags = document.querySelectorAll(".message");
	var last = msgLiTags[msgLiTags.length - 1];
	last.scrollIntoView();
}