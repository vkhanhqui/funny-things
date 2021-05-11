package com.chatapp.services;

import com.chatapp.models.User;

public interface LoginServiceInterface {

	public User handleLogin(String username, String password);
}
