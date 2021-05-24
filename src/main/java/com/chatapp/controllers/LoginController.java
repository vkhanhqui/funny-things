package com.chatapp.controllers;

import java.io.IOException;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpSession;

import com.chatapp.models.User;
import com.chatapp.services.FileServiceAbstract;
import com.chatapp.services.UserServiceInterface;
import com.chatapp.services.impl.UserService;

@WebServlet("/login")
public class LoginController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private UserServiceInterface userService = UserService.getInstance();

	public LoginController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		if (FileServiceAbstract.rootURL.isEmpty() || FileServiceAbstract.rootURL.contains("localhost")) {
			String url = request.getRequestURL().toString();
			FileServiceAbstract.rootURL = url.replaceAll("login", "files/");
			System.out.println(FileServiceAbstract.rootURL);
		}
		User user = (User) request.getSession().getAttribute("user");
		if (user != null) {
			response.sendRedirect("/chat");
		} else {
			RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/login.jsp");
			rd.forward(request, response);
		}
	}

	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String username = request.getParameter("username");
		String password = request.getParameter("password");
		User user = userService.findUser(username, password);
		String destPage = "/login";
		if (user != null) {
			HttpSession httpSession = request.getSession();
			httpSession.setAttribute("user", user);
			destPage = "/chat";
		}
		response.sendRedirect(destPage);
	}

}
