<%@ page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@ taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html>
<head>
<title>Login</title>
<link rel="icon" type="image/png"
	href="<c:url value="/static/images/icon.jpg" />">

<link rel="stylesheet"
	href="<c:url value="/static/css/style-login.css" />">
</head>
<body>
	<div class="wrapper">
		<div class="form">
			<div class="title">Register/Login</div>
			<form method="post" action="#" onsubmit="return validation();">
				<div class="input_wrap">
					<label for="input_text">Username</label>
					<div class="input_field">
						<input type="text" class="input" id="input_text">
					</div>
				</div>
				<div class="input_wrap">
					<label for="input_password">Password</label>
					<div class="input_field">
						<input type="password" class="input" id="input_password">
					</div>
				</div>
				<div class="input_wrap">
					<input type="button" id="register_btn" class="btn enabled" value="Register">
				</div>
				<div class="input_wrap">
					<span class="error_msg">Incorrect username or password.
						Please try again</span> <input type="submit" id="login_btn"
						class="btn disabled" value="Login" disabled="true">
				</div>
			</form>
		</div>
	</div>
	<script type="text/javascript"
		src="<c:url value="/static/js/login.js" />"></script>
</body>
</html>