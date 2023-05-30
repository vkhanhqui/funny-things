package com.chatapp.services.impl;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.List;

import javax.servlet.http.Part;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.User;
import com.chatapp.services.FileServiceAbstract;
import com.chatapp.services.UserServiceInterface;

public class UserService implements UserServiceInterface {
	private static UserService instance = null;

	private UserDaoInterface userDaoInterface = UserDao.getInstace();

	public synchronized static UserService getInstance() {
		if (instance == null) {
			instance = new UserService();
		}
		return instance;
	}

	private UserService() {
		File uploadDir = new File(FileServiceAbstract.rootLocation.toString());
		if (!uploadDir.exists()) {
			uploadDir.mkdir();
		}
		System.out.println("Root path: " + uploadDir.getAbsolutePath());
	}

	@Override
	public void saveUser(Boolean isRegister, String username, String password, boolean gender, Part avatar) {
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
				Files.copy(defaultAvatar.toPath(), newFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
			}
			User userEntity = new User(username, password, gender, fileName);
			userDaoInterface.saveUser(isRegister, userEntity);
		} catch (IOException ex) {
			System.out.println(ex.getMessage());
			ex.printStackTrace();
		}

	}

	@Override
	public User findUser(String username, String password) {
		User user = userDaoInterface.findByUserNameAndPassword(username, password);
		return user;
	}

	@Override
	public List<User> findFriends(String username) {
		List<User> friends = userDaoInterface.findFriends(username);
		return friends;
	}

	public List<User> findFriendsByKeyWord(String username, String keyword) {
		List<User> friends = userDaoInterface.findFriendsByKeyWord(username, keyword);
		return friends;
	}

	@Override
	public List<User> getFriendsNotInConversation(String userName, String keyword, Long conversationId) {
		List<User> friends = userDaoInterface.findFriendsNotInConversation(userName, keyword, conversationId);
		return friends;
	}
}
