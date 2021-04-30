<%@ page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@ taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html>
<head>
<title>Index</title>

<link rel="stylesheet"
	href="<c:url value="/static/css/style-chatbox-ui.css" />">
</head>
<body>
	<div class="login-form">
		<form action="#">
			<label>userId</label> <input type="text" name="userId" /><br> <label>password</label>
			<input type="password" name="password" />
			<button type="button" onclick="init();">Submit</button>
		</form>
	</div>

	<div id="container" style="display: none">
		<aside>
			<header>
				<input type="text" placeholder="Search">
			</header>
			<ul>
				<c:forEach var="userId" items="${idUsers}">
					<li id=${userId } onclick="setReceiver(this);"><img
						src="<c:url value="/static/images/chat_avatar_01.jpg" />"
						alt="${userId}">
						<div>
							<h2>${userId}</h2>
							<h3 id="status-${userId}">
								<span class="status orange"></span>
								offline
							</h3>
						</div></li>
				</c:forEach>
			</ul>
		</aside>
		<main>
			<header>
				<img src="<c:url value="/static/images/chat_avatar_01.jpg" />"
					alt="">
				<div>
					<br>
					<h2 id="receiver">Vincent Porter</h2>
				</div>
			</header>
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
		src="<c:url value="/static/js/chatbox-ui.js" />"></script>
</body>
</html>