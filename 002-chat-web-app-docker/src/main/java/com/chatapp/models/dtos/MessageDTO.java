package com.chatapp.models.dtos;

import java.util.HashSet;
import java.util.Set;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class MessageDTO {

	@JsonProperty("username")
	private String username;
	@JsonProperty("avatar")
	private String avatar;
	@JsonProperty("message")
	private String message;
	@JsonProperty("type")
	private String type;
	@JsonProperty("receiver")
	private String receiver;
	@JsonProperty("groupId")
	private Long groupId;
	@JsonProperty("onlineList")
	private Set<String> onlineList = new HashSet<String>();

	public MessageDTO() {
	}

	@JsonCreator
	public MessageDTO(@JsonProperty("username") String username, @JsonProperty("message") String message,
			@JsonProperty("type") String type, @JsonProperty("receiver") String receiver,
			@JsonProperty("groupId") Long groupId) {
		this.username = username;
		this.message = message;
		this.type = type;
		this.receiver = receiver;
		this.groupId = groupId;
	}

	public String getAvatar() {
		return avatar;
	}

	public void setAvatar(String avatar) {
		this.avatar = avatar;
	}

	public Long getGroupId() {
		return groupId;
	}

	public void setGroupId(Long groupId) {
		this.groupId = groupId;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	public String getType() {
		return type;
	}

	public void setType(String type) {
		this.type = type;
	}

	public String getReceiver() {
		return receiver;
	}

	public void setReceiver(String receiver) {
		this.receiver = receiver;
	}

	public Set<String> getOnlineList() {
		return onlineList;
	}

	public void setOnlineList(Set<String> onlineList) {
		this.onlineList = onlineList;
	}
}