package com.chatapp.models;

import java.util.Objects;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public final class Message {

	@JsonProperty("email")
	private final String email;
	@JsonProperty("message")
	private final String message;
	@JsonProperty("receiver")
	private final String receiver;

	@JsonCreator
	public Message(@JsonProperty("email") final String email, @JsonProperty("message") final String message,
			@JsonProperty("receiver") String receiver) {
		Objects.requireNonNull(email);
		Objects.requireNonNull(message);

		this.email = email;
		this.message = message;
		this.receiver = receiver;
	}

	public String getEmail() {
		return this.email;
	}

	public String getMessage() {
		return this.message;
	}

	public String getReceiver() {
		return receiver;
	}
}
