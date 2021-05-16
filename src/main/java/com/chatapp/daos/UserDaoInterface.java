package com.chatapp.daos;

import java.util.List;

import com.chatapp.models.User;

public interface UserDaoInterface extends GenericDaoInterface<User> {
	User findByUserNameAndPassword(String userName, String password);

	void saveUser(Boolean isRegister, User user);

	List<User> findFriends(String userName);
}
