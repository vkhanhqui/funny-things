package com.chatapp.models;

public class Message {

	private String username;
	private String message;
	private String avatar;
	private String type;
	private String receiver;
	private Long groupId;

	public Message() {

	}

	public Message(String username, String message, String type, String receiver, Long groupId) {
		this.username = username;
		this.message = message;
		this.type = type;
		this.receiver = receiver;
		this.groupId = groupId;
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

	public String getAvatar() {
		return avatar;
	}

	public void setAvatar(String avatar) {
		this.avatar = avatar;
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
}
