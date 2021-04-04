<%@ page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@ taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html>

<head>
<meta charset="UTF-8">
<title>Chat - Websockets</title>
<link rel="stylesheet" href="<c:url value="/static/css/style.css" />">
</head>

<body>

	<div class="login-form">
		<form action="#">
			<label>email</label> <input type="text" name="email" /><br> <label>password</label>
			<input type="password" name="password" />
			<button type="button" onclick="init();">Submit</button>
		</form>
	</div>

	<div id="main" class="main" style="display: none">
		<div id="scrolling-messages" class="scrolling-messages"></div>
		<div class="message-label">
			<span>Enter message:</span>
		</div>
		<div class="message-section">
			<div>
				<textarea id="message"></textarea>
			</div>
			<div style="float: right">
				<input type="button" value="submit" onclick="sendMessage();"
					class="button" />
			</div>
		</div>
	</div>

	<script type="text/javascript"
		src="<c:url value="/static/js/chat.js" />"></script>
</body>
</html>