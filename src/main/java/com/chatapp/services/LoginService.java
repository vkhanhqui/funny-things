package com.chatapp.services;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.User;

public class LoginService {
	private static LoginService instance = null;
	
	private UserDaoInterface userDaoInterface = UserDao.getInstace();

	public synchronized static LoginService getInstance() {
		if (instance == null) {
			instance = new LoginService();
		}
		return instance;
	}

	private LoginService() {
	}

	public User handleLogin(String username, String password) {
		User user = userDaoInterface.findByUserNameAndPassword(username, password);
		return user;
	}

}
