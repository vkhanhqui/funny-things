package com.chatapp.mapper;

import java.sql.ResultSet;
import java.sql.SQLException;

import com.chatapp.models.User;

public class UserMapper implements RowMapper<User> {

	@Override
	public User mapRow(ResultSet resultSet) {
		try {
			User user = new User();
			user.setUsername(resultSet.getString("username"));
			user.setGender(resultSet.getBoolean("gender"));
			user.setAvatar(resultSet.getString("avatar"));
			return user;
		} catch (SQLException e) {
			return null;
		}
	}
}
