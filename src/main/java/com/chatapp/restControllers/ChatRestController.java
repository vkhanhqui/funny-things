package com.chatapp.restControllers;

import java.io.IOException;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.List;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.Message;
import com.chatapp.services.ChatService;
import com.fasterxml.jackson.databind.ObjectMapper;

@WebServlet("/chat-rest-controller")
public class ChatRestController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	public ChatRestController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {		
		List<Message> messages = new ArrayList<>();
		//current user is a1
		messages.add(new Message("a1", "hello", "a2", ChatService.onlineList));
		//clicked user
		String userId = request.getParameter("userId");
		for (int i = 0; i < 10; i++) {
			messages.add(new Message(userId, "hello", "a1", ChatService.onlineList));
		}
		//response to json
		ObjectMapper objectMapper = new ObjectMapper();
		String json = objectMapper.writeValueAsString(messages);

		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");

		PrintWriter printWriter = response.getWriter();
		printWriter.print(json);
		printWriter.flush();
	}

}
