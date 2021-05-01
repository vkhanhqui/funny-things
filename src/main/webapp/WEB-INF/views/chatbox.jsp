<%@ page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@ taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html>
<head>
<title>Chat</title>
<link rel="icon" type="image/png"
	href="<c:url value="/static/images/icon.jpg" />">

<link rel="stylesheet"
	href="<c:url value="/static/css/style-chatbox.css" />">
</head>
<body onload="init()">
	<p id="username" style="display: none">${user.username}</p>

	<div id="container">
		<aside>
			<header>
				<input type="text" placeholder="Search">
			</header>
			<ul>
				<c:forEach var="friend" items="${friends}">
					<li id=${friend } onclick="setReceiver(this);"><img
						src="<c:url value="/static/images/chat_avatar_01.jpg" />"
						alt="${friend}">
						<div>
							<h2>${friend}</h2>
							<h3 id="status-${friend}">
								<span class="status orange"></span> offline
							</h3>
						</div></li>
				</c:forEach>
			</ul>
		</aside>
		<main>
			<header id="receiver"> </header>
			<ul id="chat">

			</ul>
			<footer>
				<textarea placeholder="Type your message" id="message"></textarea>
				<img
					src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/1940306/ico_picture.png"
					alt=""> <img
					src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/1940306/ico_file.png"
					alt=""> <a href="#" onclick="sendMessage();">Send</a>
			</footer>
		</main>
	</div>

	<script type="text/javascript"
		src="<c:url value="/static/js/chatbox.js" />"></script>
</body>
</html>