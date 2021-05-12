package com.chatapp.services;

import javax.servlet.http.Part;

public interface UserServiceInterface {

	public void saveUser(Boolean isRegister, String username, String password, boolean gender, Part avatar);
}
