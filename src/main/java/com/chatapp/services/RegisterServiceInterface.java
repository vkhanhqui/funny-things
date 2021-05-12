package com.chatapp.services;

import javax.servlet.http.Part;

public interface RegisterServiceInterface {

	public void handleRegister(String username, String password, boolean gender, Part avatar);
}
