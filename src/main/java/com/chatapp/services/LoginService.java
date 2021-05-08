package com.chatapp.services;

import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

import com.chatapp.daos.DBConnection;
import com.chatapp.models.User;

public class LoginService {
	private static LoginService instance = null;

	public synchronized static LoginService getInstance() {
		if (instance == null) {
			instance = new LoginService();
		}
		return instance;
	}

	private LoginService() {
	}

	public User handleLogin(String username, String password) {
		User user = null;
		if (DBConnection.Open()) {
			String sql = "select username, gender, avatar " + "from testing " + "where username=? and password=?";
			try {
				PreparedStatement ps = DBConnection.connection.prepareStatement(sql);
				ps.setString(1, username);
				ps.setString(2, password);
				ResultSet rs = ps.executeQuery();
				if (rs.next()) {
					boolean gender = rs.getBoolean(2);
					String avatar = rs.getString(3);
					user = new User(username, gender, avatar);
				}
			} catch (SQLException ex) {

			} finally {
				DBConnection.Close();
			}
		}
		return user;
	}

}
