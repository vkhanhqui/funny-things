package com.chatapp.controllers;

import java.io.IOException;
import java.util.List;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.daos.impl.FriendDao;
import com.chatapp.models.Friend;
import com.chatapp.models.User;
import com.chatapp.services.UserServiceInterface;
import com.chatapp.services.impl.UserService;

@WebServlet("/chat")
public class ChatController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private UserServiceInterface userService = UserService.getInstance();
	private FriendDao friendDao = new FriendDao();
	
	public ChatController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		User currentUser = (User) request.getSession().getAttribute("user");
		List<User> friends = userService.findFriends(currentUser.getUsername());

		request.setAttribute("friends", friends);
		request.setAttribute("user", currentUser);

		RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/chatbox.jsp");
		rd.forward(request, response);
	}
	
	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {

		String sender = request.getParameter("sender");
		String receiver = request.getParameter("receiver");
		boolean status = Boolean.parseBoolean(request.getParameter("status"));
		boolean isAccept = Boolean.parseBoolean(request.getParameter("isAccept"));
		friendDao.saveFriend(isAccept, new Friend(sender, receiver, sender, status));

		response.sendRedirect("/chat");
	}
}
