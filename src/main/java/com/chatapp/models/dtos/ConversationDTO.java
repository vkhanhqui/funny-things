package com.chatapp.models.dtos;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class ConversationDTO {
	@JsonProperty("name")
	private String name;

	@JsonProperty("users")
	private List<UserDTO> users;

	@JsonCreator
	public ConversationDTO(@JsonProperty("name") String name, @JsonProperty("users") List<UserDTO> users) {
		super();
		this.name = name;
		this.users = users;
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	public List<UserDTO> getUsers() {
		return users;
	}

	public void setUsers(List<UserDTO> users) {
		this.users = users;
	}
}
