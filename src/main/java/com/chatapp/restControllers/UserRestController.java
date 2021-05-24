package com.chatapp.restControllers;

import java.io.IOException;
import java.io.PrintWriter;
import java.util.List;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.models.User;
import com.chatapp.services.UserServiceInterface;
import com.chatapp.services.impl.ChatService;
import com.chatapp.services.impl.UserService;
import com.fasterxml.jackson.databind.ObjectMapper;

@WebServlet("/users-rest-controller")
public class UserRestController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private UserServiceInterface userServiceInterface = UserService.getInstance();
	private ChatService chatService = ChatService.getInstance();

	public UserRestController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String keyWord = request.getParameter("keyword");
		String userName = request.getParameter("username");
		List<User> listUsers;
		if (keyWord.isEmpty()) {
			listUsers = userServiceInterface.findFriends(userName);
		} else {
			listUsers = userServiceInterface.findFriendsByKeyWord(userName, keyWord);
		}
		for (User user : listUsers) {
			user.setOnline(chatService.isUserOnline(user.getUsername()));
		}
		ObjectMapper objectMapper = new ObjectMapper();
		String json = objectMapper.writeValueAsString(listUsers);

		response.setContentType("application/json");
		response.setCharacterEncoding("UTF-8");

		PrintWriter printWriter = response.getWriter();
		printWriter.print(json);
		printWriter.flush();
	}
}
