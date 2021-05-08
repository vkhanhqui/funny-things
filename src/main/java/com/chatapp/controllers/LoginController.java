package com.chatapp.controllers;

import java.io.IOException;

import javax.servlet.RequestDispatcher;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpSession;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.User;
import com.chatapp.services.LoginService;

@WebServlet("/login")
public class LoginController extends HttpServlet {
	private static final long serialVersionUID = 1L;
	private static LoginService loginService = LoginService.getInstance();
	
	private UserDaoInterface userDaoInterface = UserDao.getInstace();

	public LoginController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
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
		
		if(userDaoInterface.findByUserNameAndPassword(username, password) != null) {
			System.out.println("okakakakakkakakalla");
		}
		
		User user = loginService.handleLogin(username, password);
		String destPage = "/login";
		if (user != null) {
			HttpSession httpSession = request.getSession();
			httpSession.setAttribute("user", user);
			destPage = "/chat";
		}
		response.sendRedirect(destPage);
	}

}
