package com.chatapp.models.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

public class UserDTO {
	private String username;
	private boolean isAdmin;

	public UserDTO(@JsonProperty("username") String username, @JsonProperty("isAdmin") boolean isAdmin) {
		this.username = username;
		this.isAdmin = isAdmin;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}

	public boolean isAdmin() {
		return isAdmin;
	}

	public void setAdmin(boolean isAdmin) {
		this.isAdmin = isAdmin;
	}

}
