package com.chatapp.mappers.impl;

import java.sql.ResultSet;
import java.sql.SQLException;

import com.chatapp.mappers.RowMapperInterface;
import com.chatapp.models.User;

public class FriendMapper implements RowMapperInterface<User> {

	@Override
	public User mapRow(ResultSet resultSet) {
		try {
			User user = new User();
			user.setUsername(resultSet.getString("u1").trim());
			user.setAvatar(resultSet.getString("u1_avt").trim());
			user.setUsername(resultSet.getString("u2").trim());
			user.setAvatar(resultSet.getString("u2_avt").trim());
			
			return user;
		} catch (SQLException e) {
			return null;
		}
	}
}
