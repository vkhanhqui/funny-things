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

import com.chatapp.services.UserServiceInterface;
import com.chatapp.services.impl.UserService;

@WebServlet("/users/*")
@MultipartConfig
public class UserController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private UserServiceInterface userService = UserService.getInstance();

	public UserController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String status = request.getPathInfo();
		if (status.equals("/register") || status.equals("/update")) {
			String title = "Update User";
			String description = "Update your information";
			String btnSubmit = "Update";
			String btnGoBack = "/chat";

			if (status.equals("/register")) {
				title = "Register User";
				description = "Enter your information";
				btnSubmit = "Register";
				btnGoBack = "/login";
			}

			request.setAttribute("title", title);
			request.setAttribute("description", description);
			request.setAttribute("status", status);
			request.setAttribute("btnSubmit", btnSubmit);
			request.setAttribute("btnGoBack", btnGoBack);

			RequestDispatcher rd = request.getRequestDispatcher("/WEB-INF/views/user-form.jsp");
			rd.forward(request, response);
		} else if (status.equals("/logout")) {
			request.getSession().invalidate();
			response.sendRedirect("/login");
		} else {
			response.sendRedirect("/login");
		}
	}

	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String username = request.getParameter("username");
		String password = request.getParameter("password");
		String gender = request.getParameter("gender");
		Part avatar = request.getPart("avatar");

		String path = request.getPathInfo();
		if (path.endsWith("register")) {
			userService.saveUser(true, username, password, Boolean.valueOf(gender), avatar);
			response.sendRedirect("/login");
		} else if (path.endsWith("update")) {
			userService.saveUser(false, username, password, Boolean.valueOf(gender), avatar);
			response.sendRedirect("/users/update");
		} else {
			response.sendRedirect("/chat");
		}
	}
}