package com.chatapp.models;

import java.util.HashSet;
import java.util.Objects;
import java.util.Set;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class Message {

	@JsonProperty("username")
	private String username;
	@JsonProperty("message")
	private String message;
	@JsonProperty("receiver")
	private String receiver;
	@JsonProperty("onlineList")
	private Set<String> onlineList = new HashSet<String>();

	@JsonCreator
	public Message(@JsonProperty("username") String username, @JsonProperty("message") String message,
			@JsonProperty("receiver") String receiver, @JsonProperty("onlineList") Set<String> onlineList) {
		Objects.requireNonNull(username);
		Objects.requireNonNull(message);
		Objects.requireNonNull(receiver);

		this.username = username;
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

	public String getusername() {
		return username;
	}

	public void setusername(String username) {
		this.username = username;
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
