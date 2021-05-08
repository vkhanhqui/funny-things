package com.chatapp.daos.impl;

import java.util.List;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.mapper.UserMapper;
import com.chatapp.models.User;

public class UserDao extends AbstractDao<User> implements UserDaoInterface {
	
	private static UserDao instance = null;
	
	private UserDao() {
		
	}
	
	public synchronized static UserDao getInstace() {
		if(instance == null) {
			instance = new UserDao();
		}
		return instance;
	}
	
	@Override
	public User findByUserNameAndPassword(String userName, String password) {
		StringBuilder sql = new StringBuilder("select username, gender, avatar");
		sql.append(" from testing where username=? and password=?");
		List<User> users = query(sql.toString(), new UserMapper(), userName, password);
		return users.isEmpty() ? null : users.get(0);
	}

}
