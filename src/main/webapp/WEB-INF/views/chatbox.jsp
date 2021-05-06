<%@ page language="java" contentType="text/html; charset=UTF-8"
	pageEncoding="UTF-8"%>
<%@ taglib uri="http://java.sun.com/jsp/jstl/core" prefix="c"%>
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link
	href="//maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"
	rel="stylesheet">
<link rel="stylesheet" href="<c:url value="/static/css/style.css" />">

<link rel="icon" type="image/png"
	href="<c:url value="/static/images/icon.jpg" />">
<title>Chat</title>
</head>
<body>
	<p id="username" style="display: none">${user.username}</p>
	<div class="container">
		<div class="conversation-container">
			<div class="left-side">
				<h2>Chats</h2>
				<div class="tab-control">
					<i class="fa fa-comment active"></i> <i class="fa fa-comments"></i>
				</div>
				<div class="list-user-search">
					<input type="text" class="txt-input" placeholder="Search...">
				</div>
				<div class="list-user">
					<ul>
						<c:forEach var="friend" items="${friends}">
							<li id=${friend } onclick="setReceiver(this);">
								<div class="user-contain">
									<div class="user-img">
										<img src="<c:url value="/static/images/user-male.jpg" />"
											alt="Image of user">
										<div id="status-${friend}" class="user-img-dot"></div>
									</div>
									<div class="user-info">
										<span class="user-name">${friend}</span> <span
											class="user-last-message">Hello!</span>
									</div>
								</div>
							</li>
						</c:forEach>
					</ul>
				</div>
			</div>
			<div class="right-side active" id="receiver"></div>
		</div>
	</div>


	<script type="text/javascript"
		src="<c:url value="/static/js/chatbox.js" />" charset="utf-8"></script>
</body>
</html>