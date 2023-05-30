package com.chatapp.restControllers;

import java.io.IOException;
import java.io.PrintWriter;
import java.util.List;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.services.MessageServiceInterface;
import com.chatapp.services.impl.MessageService;
import com.fasterxml.jackson.databind.ObjectMapper;

@WebServlet("/chat-rest-controller")
public class ChatRestController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private MessageServiceInterface messageServiceInterface = MessageService.getInstance();

	public ChatRestController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String sender = request.getParameter("sender");
		String receiver = request.getParameter("receiver");
		List<MessageDTO> messages = messageServiceInterface.getAllMessagesBySenderAndReceiver(sender, receiver);

		ObjectMapper objectMapper = new ObjectMapper();
		String json = objectMapper.writeValueAsString(messages);

		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");

		PrintWriter printWriter = response.getWriter();
		printWriter.print(json);
		printWriter.flush();
	}
	
}
