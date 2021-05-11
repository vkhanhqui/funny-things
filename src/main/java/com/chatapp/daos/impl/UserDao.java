package com.chatapp.daos.impl;

import java.util.List;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.mappers.impl.UserMapper;
import com.chatapp.models.User;

public class UserDao extends GenericDao<User> implements UserDaoInterface {

	private static UserDao instance = null;

	private UserDao() {

	}

	public synchronized static UserDao getInstace() {
		if (instance == null) {
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

	@Override
	public void saveUser(String username, String password, Boolean gender, String avatar) {
		StringBuilder sql = new StringBuilder("insert into testing values(?,?,?,?)");
		insert(sql.toString(), username, password, gender, avatar);
	}

}
