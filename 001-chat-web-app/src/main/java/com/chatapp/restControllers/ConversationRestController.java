package com.chatapp.restControllers;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.List;
import java.util.stream.Collectors;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.models.dtos.UserDTO;
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

	@Override
	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String username = request.getParameter("username");
		String usersConversationId = request.getParameter("usersConversationId");
		String messagesConversationId = request.getParameter("messagesConversationId");
		String conversationKeyword = request.getParameter("conversationKeyword");
		String json = "Must have username or conversation id as request param";

		ObjectMapper objectMapper = new ObjectMapper();

		if (conversationKeyword != null && !conversationKeyword.isEmpty() && username != null && !username.isEmpty()) {

			List<ConversationDTO> conversationDTOs = conversationServiceInterface
					.getConversationsOfUserByKeyword(username, conversationKeyword);
			for (ConversationDTO c : conversationDTOs) {
				List<UserDTO> userDTOs = conversationServiceInterface.getAllUsersByConversationId(c.getId()).stream()
						.filter(u -> u.isAdmin()).collect(Collectors.toList());
				c.setUsers(userDTOs);
			}
			json = objectMapper.writeValueAsString(conversationDTOs);

		} else if (username != null && !username.isEmpty()) {

			List<ConversationDTO> conversationDTOs = conversationServiceInterface
					.getAllConversationsByUsername(username);
			for (ConversationDTO conversationDTO : conversationDTOs) {
				conversationDTO
						.setUsers(conversationServiceInterface.getAllUsersByConversationId(conversationDTO.getId()));
			}
			json = objectMapper.writeValueAsString(conversationDTOs);

		} else if (usersConversationId != null && !usersConversationId.isEmpty()) {

			Long id = Long.parseLong(usersConversationId);
			List<UserDTO> userDTOs = conversationServiceInterface.getAllUsersByConversationId(id);

			json = objectMapper.writeValueAsString(userDTOs);

		} else if (messagesConversationId != null && !messagesConversationId.isEmpty()) {

			Long id = Long.parseLong(messagesConversationId);
			List<MessageDTO> messageDTOs = conversationServiceInterface.getAllMessagesByConversationId(id);

			json = objectMapper.writeValueAsString(messageDTOs);

		}

		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");
		PrintWriter printWriter = response.getWriter();
		printWriter.print(json);
		printWriter.flush();
	}

	@Override
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

	@Override
	protected void doDelete(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String username = request.getParameter("username");
		String conversationId = request.getParameter("conversationId");
		String json = "Must have username or conversation id as request param";
		if (conversationId != null && !conversationId.isEmpty()) {
			Long id = Long.parseLong(conversationId);
			if (username != null && !username.isEmpty()) {
				conversationServiceInterface.deleteUserFromConversation(id, username);
				json = "delete User by Id " + username + " From Conversation by Id " + id + " successfully";
			} else {
				conversationServiceInterface.deleteConversationById(id);
				json = "delete Conversation By Id " + id + " successfully";
			}
		}
		ObjectMapper objectMapper = new ObjectMapper();
		json = objectMapper.writeValueAsString(json);
		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");
		PrintWriter printWriter = response.getWriter();
		printWriter.print(json);
		printWriter.flush();
	}
}
