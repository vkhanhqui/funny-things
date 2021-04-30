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

@WebServlet("/chat-controller")
public class ChatController extends HttpServlet {
	private static final long serialVersionUID = 1L;	

	public ChatController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		List<String> idUsers = new ArrayList<>();
		for(int i=1; i<=20; i++) {
			idUsers.add("a"+i);
		}
		request.setAttribute("idUsers", idUsers);
		
		RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/chatbox-ui.jsp");
		rd.forward(request, response);
	}

}
