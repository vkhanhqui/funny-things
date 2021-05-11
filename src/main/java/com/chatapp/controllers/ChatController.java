package com.chatapp.controllers;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.User;

@WebServlet("/chat")
public class ChatController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	public ChatController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		List<User> friends = new ArrayList<>();
		for (int i = 1; i <= 11; i++) {
			User newUser = new User();
			newUser.setUsername("a" + i);
			newUser.setAvatar("a" + i+".jpg");
			friends.add(newUser);
		}
		request.setAttribute("friends", friends);
		User user = (User) request.getSession().getAttribute("user");
		request.setAttribute("user", user);

		RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/chatbox.jsp");
		rd.forward(request, response);
	}

}
