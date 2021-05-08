package com.chatapp.services;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.sql.PreparedStatement;
import java.sql.SQLException;

import javax.servlet.http.Part;

import com.chatapp.daos.DBConnection;

public class RegisterService {
	private static RegisterService instance = null;
	public static Path rootLocation = Paths.get("archive");

	public synchronized static RegisterService getInstance() {
		if (instance == null) {
			instance = new RegisterService();
		}
		return instance;
	}

	private RegisterService() {
		File uploadDir = new File(rootLocation.toString());
		if (!uploadDir.exists()) {
			uploadDir.mkdir();
		}
		System.out.println("Root path: " + uploadDir.getAbsolutePath());
	}

	public void handleRegister(String username, String password, boolean gender, Part avatar) {
		if (DBConnection.Open()) {
			try {
				File privateDir = new File(rootLocation.toString() + "/" + username);
				privateDir.mkdir();
				String origin = avatar.getSubmittedFileName();
				String fileName = "";
				if (!origin.isEmpty()) {
					String tail = origin.substring(origin.lastIndexOf("."), origin.length());
					fileName = username + tail;
					avatar.write(privateDir.getAbsolutePath() + File.separator + fileName);
				} else {
					File defaultAvatar = new File(rootLocation.toString() + "/default/user-male.jpg");
					if (gender == false) {
						defaultAvatar = new File(rootLocation.toString() + "/default/user-female.jpg");
					}
					fileName = username + ".jpg";
					File newFile = new File(privateDir.toString() + "/" + fileName);
					Files.copy(defaultAvatar.toPath(), newFile.toPath());
				}
				String sql = "insert into testing values(?,?,?,?)";
				PreparedStatement ps = DBConnection.connection.prepareStatement(sql);
				ps.setString(1, username);
				ps.setString(2, password);
				ps.setBoolean(3, gender);
				ps.setString(4, fileName);
				ps.executeUpdate();
			} catch (IOException | SQLException ex) {
				System.out.println(ex.getMessage());
				ex.printStackTrace();
			} finally {
				DBConnection.Close();
			}
		}
	}
}
