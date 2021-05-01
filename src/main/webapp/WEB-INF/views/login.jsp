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
			<div class="title">Login</div>
			<form method="post" action="<c:url value="/login" />">

				<div class="input_wrap">
					<label for="input_text">Username</label>
					<div class="input_field">
						<input type="text" class="input" onchange="validateUsername(this)">
					</div>
				</div>

				<div class="input_wrap">
					<label for="input_password">Password</label>
					<div class="input_field">
						<input type="password" class="input"
							onchange="validatePassword(this)">
					</div>
				</div>

				<div class="input_wrap">
					<input type="submit" class="btn disabled" value="Login">
				</div>

				<div class="input_wrap">
					<input type="button" class="btn enabled" value="Or Register"
						onclick="changeLoginForm(true)">
				</div>

			</form>
		</div>

		<div class="form" style="display: none">
			<div class="title">Register</div>
			<form method="post" action="<c:url value="/register" />">
				<div class="input_wrap">
					<label for="input_text">Username</label>
					<div class="input_field">
						<input type="text" class="input" onchange="validateUsername(this)">
					</div>
				</div>
				<div class="input_wrap">
					<label for="input_password">Password</label>
					<div class="input_field">
						<input type="password" class="input"
							onchange="validatePassword(this)">
					</div>
				</div>
				<div class="input_wrap">
					<label for="input_password">Confirm Password</label>
					<div class="input_field">
						<input type="password" class="input"
							onchange="validateConfirmPassword(this)">
					</div>
				</div>

				<div class="input_wrap">
					<label for="genders">Gender </label>
					<div class="input_field">
						<select class="input" onclick="loadDefaultImage(this)">
							<option value="male" class="input">Male</option>
							<option value="female" class="input">Female</option>
						</select>
					</div>
				</div>

				<div class="input_wrap">
					<label for="file">Upload Your Avatar</label> <input type="file"
						accept="image/*" name="image" onchange="loadImage(event)">
					<div class="input_field">
						<img id="display-image"
							src="<c:url value="/static/images/default-user-male.png" />" />
					</div>
				</div>

				<div class="input_wrap">
					<input type="submit" class="btn disabled" value="Register">
				</div>

				<div class="input_wrap">
					<input type="button" class="btn enabled" value="Or Login"
						onclick="changeLoginForm(false)">
				</div>

			</form>
		</div>
	</div>
	<script type="text/javascript"
		src="<c:url value="/static/js/login.js" />"></script>
</body>
</html>