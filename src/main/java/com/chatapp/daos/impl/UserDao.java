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
	public void saveUser(Boolean isRegister, User user) {
		String username = user.getUsername();
		String password = user.getPassword();
		Boolean gender = user.isGender();
		String avatar = user.getAvatar();
		StringBuilder sql = new StringBuilder("insert into testing values(?,?,?,?)");
		if (isRegister) {
			save(sql.toString(), username, password, gender, avatar);
		} else {
			sql = new StringBuilder("update testing set password=?, gender=?, avatar=? where username=?");
			save(sql.toString(), password, gender, avatar, username);
		}
	}

}
