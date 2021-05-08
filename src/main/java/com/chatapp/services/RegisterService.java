package com.chatapp.services;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import javax.servlet.http.Part;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;

public class RegisterService {
	private static RegisterService instance = null;
	public static Path rootLocation = Paths.get("archive");

	private UserDaoInterface userDaoInterface = UserDao.getInstace();

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
			userDaoInterface.saveUser(username, password, gender, fileName);
		} catch (IOException ex) {
			System.out.println(ex.getMessage());
			ex.printStackTrace();
		}

	}
}
