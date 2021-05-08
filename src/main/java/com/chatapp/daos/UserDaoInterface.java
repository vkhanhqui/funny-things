package com.chatapp.daos;

import com.chatapp.models.User;

public interface UserDaoInterface extends GenericDaoInterface<User> {
	User findByUserNameAndPassword(String userName, String password);
}
