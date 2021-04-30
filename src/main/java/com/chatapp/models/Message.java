package com.chatapp.models;

import java.util.HashSet;
import java.util.Objects;
import java.util.Set;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class Message {

	@JsonProperty("userId")
	private String userId;
	@JsonProperty("message")
	private String message;
	@JsonProperty("receiver")
	private String receiver;
	@JsonProperty("onlineList")
	private Set<String> onlineList = new HashSet<String>();

	@JsonCreator
	public Message(@JsonProperty("userId") String userId, @JsonProperty("message") String message,
			@JsonProperty("receiver") String receiver, @JsonProperty("onlineList") Set<String> onlineList) {
		Objects.requireNonNull(userId);
		Objects.requireNonNull(message);
		Objects.requireNonNull(receiver);

		this.userId = userId;
		this.message = message;
		this.receiver = receiver;
		this.onlineList = onlineList;
	}

	public String getMessage() {
		return this.message;
	}

	public String getReceiver() {
		return receiver;
	}

	public String getUserId() {
		return userId;
	}

	public void setUserId(String userId) {
		this.userId = userId;
	}

	public void setMessage(String message) {
		this.message = message;
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
