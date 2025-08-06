package com.chatapp.models.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

public class UserDTO {
	private String username;
	private String avatar;
	private boolean isAdmin;
	
	public UserDTO() {
	}

	public UserDTO(@JsonProperty("username") String username, @JsonProperty("avatar") String avatar,
			@JsonProperty("isAdmin") boolean isAdmin) {
		this.username = username;
		this.isAdmin = isAdmin;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}

	public String getAvatar() {
		return avatar;
	}

	public void setAvatar(String avatar) {
		this.avatar = avatar;
	}

	public boolean isAdmin() {
		return isAdmin;
	}

	public void setAdmin(boolean isAdmin) {
		this.isAdmin = isAdmin;
	}

}
