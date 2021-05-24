package com.chatapp.restControllers;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.PrintWriter;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.services.ConversationServiceInterface;
import com.chatapp.services.impl.ConversationService;
import com.fasterxml.jackson.databind.ObjectMapper;

@WebServlet("/conversations-rest-controller")
public class ConversationRestController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private ConversationServiceInterface conversationServiceInterface = ConversationService.getInstance();

	public ConversationRestController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		
	}

	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");
		PrintWriter printWriter = response.getWriter();
		String json = "";

		StringBuilder requestBody = new StringBuilder();
		String line = null;
		try {
			BufferedReader reader = request.getReader();
			while ((line = reader.readLine()) != null) {
				requestBody.append(line);
			}
		} catch (IOException ex) {
			json = ex.getMessage();
			printWriter.print(json);
			printWriter.flush();
		}

		ObjectMapper objectMapper = new ObjectMapper();
		ConversationDTO conversationDTO = objectMapper.readValue(requestBody.toString(), ConversationDTO.class);
		conversationServiceInterface.saveConversation(conversationDTO);
		json = objectMapper.writeValueAsString(conversationDTO);

		printWriter.print(json);
		printWriter.flush();
	}
}
