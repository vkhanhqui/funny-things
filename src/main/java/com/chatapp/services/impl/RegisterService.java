package com.chatapp.services.impl;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;

import javax.servlet.http.Part;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.services.FileServiceAbstract;
import com.chatapp.services.RegisterServiceInterface;

public class RegisterService implements RegisterServiceInterface{
	private static RegisterService instance = null;

	private UserDaoInterface userDaoInterface = UserDao.getInstace();

	public synchronized static RegisterService getInstance() {
		if (instance == null) {
			instance = new RegisterService();
		}
		return instance;
	}

	private RegisterService() {
		File uploadDir = new File(FileServiceAbstract.rootLocation.toString());
		if (!uploadDir.exists()) {
			uploadDir.mkdir();
		}
		System.out.println("Root path: " + uploadDir.getAbsolutePath());
	}

	@Override
	public void handleRegister(String username, String password, boolean gender, Part avatar) {
		try {
			File privateDir = new File(FileServiceAbstract.rootLocation.toString() + "/" + username);
			privateDir.mkdir();
			String origin = avatar.getSubmittedFileName();
			String fileName = "";
			if (!origin.isEmpty()) {
				String tail = origin.substring(origin.lastIndexOf("."), origin.length());
				fileName = username + tail;
				avatar.write(privateDir.getAbsolutePath() + File.separator + fileName);
			} else {
				File defaultAvatar = new File(FileServiceAbstract.rootLocation.toString() + "/default/user-male.jpg");
				if (gender == false) {
					defaultAvatar = new File(FileServiceAbstract.rootLocation.toString() + "/default/user-female.jpg");
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
