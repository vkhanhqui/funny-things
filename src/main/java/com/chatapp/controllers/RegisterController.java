package com.chatapp.controllers;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.annotation.MultipartConfig;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.Part;

import com.chatapp.services.RegisterServiceInterface;
import com.chatapp.services.impl.RegisterService;

@WebServlet("/register")
@MultipartConfig
public class RegisterController extends HttpServlet {
	private static final long serialVersionUID = 1L;
	private RegisterServiceInterface registerService = RegisterService.getInstance();

	public RegisterController() {
		super();
	}

	protected void doPost(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String username = request.getParameter("username");
		String password = request.getParameter("password");
		String gender = request.getParameter("gender");
		Part avatar = request.getPart("avatar");
		registerService.handleRegister(username, password, Boolean.valueOf(gender), avatar);

		response.sendRedirect("/login");
	}
}
