<%@page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="<c:url value="/static/css/style.css" />">
<link rel="icon" type="image/png"
	href="<c:url value="/static/images/icon.jpg" />">

<title>Edit Conversation</title>
</head>
<body>
	<div class="container">
		<div class="form-container">
			<div class="tab-control">
				<h3 class="tab-control-btn register">Edit Conversation</h3>
			</div>
			<div class="register-form form active">
				<form action="<c:url value="/conversation" />"
					enctype="multipart/form-data" method="POST">

					<input type="hidden" class="txt-input border" name="groupId"
						value="${conversationDTO.id}"> <input type="text"
						class="txt-input border" name="groupName"
						value="${conversationDTO.name}"> <label for="image">
						<img
						src="<c:url value="/files/group-${conversationDTO.id}/${conversationDTO.avatar}" />"
						class="image-profile" alt="">
					</label> <input type="file" accept="image/*" name="avatar" id="image"
						class="image-file">

					<button type="submit" class="btn btn-login border">Submit</button>

					<a href="<c:url value="/chat" />" class="btn btn-login border"
						style="background-color: grey; text-align: center;">Go back</a>

				</form>
			</div>
		</div>
	</div>

	<script type="text/javascript"
		src="<c:url value="/static/js/user-form.js" />" charset="utf-8"></script>
</body>
</html>