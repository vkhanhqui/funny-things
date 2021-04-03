package com.chatapp.models;

import java.util.Objects;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public final class Message {

	@JsonProperty("email")
	private final String email;
	@JsonProperty
	private final String message;

	@JsonCreator
	public Message(@JsonProperty("email") final String email, @JsonProperty("message") final String message) {
		Objects.requireNonNull(email);
		Objects.requireNonNull(message);

		this.email = email;
		this.message = message;
	}

	public String getemail() {
		return this.email;
	}

	public String getMessage() {
		return this.message;
	}
}
