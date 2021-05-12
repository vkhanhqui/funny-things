package com.chatapp.services.impl;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.User;
import com.chatapp.services.LoginServiceInterface;

public class LoginService implements LoginServiceInterface{
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

	@Override
	public User handleLogin(String username, String password) {
		User user = userDaoInterface.findByUserNameAndPassword(username, password);
		return user;
	}

}
