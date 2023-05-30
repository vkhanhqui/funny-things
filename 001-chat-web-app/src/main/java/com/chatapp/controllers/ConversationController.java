package com.chatapp.controllers;

import java.io.IOException;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.MultipartConfig;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.Part;

import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.services.ConversationServiceInterface;
import com.chatapp.services.impl.ConversationService;

@WebServlet("/conversation")
@MultipartConfig
public class ConversationController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private ConversationServiceInterface conversationServiceInterface = ConversationService.getInstance();

	public ConversationController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String conversationId = request.getParameter("conversationId");
		String destPage = "/WEB-INF/views/group-form.jsp";
		if (conversationId != null && !conversationId.isEmpty()) {
			Long id = Long.parseLong(conversationId);
			ConversationDTO conversationDTO = conversationServiceInterface.getConversationById(id);
			request.setAttribute("conversationDTO", conversationDTO);
			RequestDispatcher rd = request.getRequestDispatcher(destPage);
			rd.forward(request, response);
		} else {
			destPage = "/chat";
			response.sendRedirect(destPage);
		}
	}

	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {

		String conversationId = request.getParameter("groupId");
		String groupName = request.getParameter("groupName").trim();
		Part avatar = request.getPart("avatar");
		Long id = Long.parseLong(conversationId);
		conversationServiceInterface.updateConversationById(id, groupName, avatar);

		response.sendRedirect("/chat");
	}
}
