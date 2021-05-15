package com.chatapp.controllers;

import java.io.IOException;
import java.util.List;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.User;

@WebServlet("/chat")
public class ChatController extends HttpServlet {
	private static final long serialVersionUID = 1L;
	
	private UserDaoInterface userDao = UserDao.getInstace();
	public ChatController() {
		super(); 
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		User currentUser = (User) request.getSession().getAttribute("user");
		List<User> friends = userDao.findFriends(currentUser.getUsername());	
		
		request.setAttribute("friends", friends);
		request.setAttribute("user", currentUser);

		RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/chatbox.jsp");
		rd.forward(request, response);
	}

}
