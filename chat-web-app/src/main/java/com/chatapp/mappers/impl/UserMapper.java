package com.chatapp.mappers.impl;

import java.sql.ResultSet;
import java.sql.SQLException;

import com.chatapp.mappers.RowMapperInterface;
import com.chatapp.models.User;

public class UserMapper implements RowMapperInterface<User> {

	@Override
	public User mapRow(ResultSet resultSet) {
		try {
			User user = new User();
			user.setUsername(resultSet.getString("username").trim());
			user.setGender(resultSet.getBoolean("gender"));
			user.setAvatar(resultSet.getString("avatar").trim());
			try {
				user.setAdmin(resultSet.getBoolean("is_admin"));
			} catch (SQLException ex) {
				return user;
			}
			return user;
		} catch (SQLException e) {
			return null;
		}
	}
}
