package com.chatapp.models;

public class User {
	private Integer id;
	private String username;
	private String password;
	private boolean gender;
	private String avatar;

	public User() {

	}

	public User(String username, String password, boolean gender, String avatar) {
		this.username = username;
		this.password = password;
		this.gender = gender;
		this.avatar = avatar;
	}

	public Integer getId() {
		return id;
	}

	public void setId(Integer id) {
		this.id = id;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}

	public boolean isGender() {
		return gender;
	}

	public void setGender(boolean gender) {
		this.gender = gender;
	}

	public String getAvatar() {
		return avatar;
	}

	public void setAvatar(String avatar) {
		this.avatar = avatar;
	}

	public String getPassword() {
		return password;
	}

	public void setPassword(String password) {
		this.password = password;
	}
}
