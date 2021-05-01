package com.chatapp.models;

import com.fasterxml.jackson.databind.ObjectMapper;

public final class Constants {

	public static final String INSTANTIATION_NOT_ALLOWED = "Instantiation not allowed";
	public static final String USERNAME_KEY = "username";
	public static final String MESSAGE_KEY = "message";
	public static final ObjectMapper MAPPER = new ObjectMapper();

	private Constants() {
		throw new IllegalStateException(INSTANTIATION_NOT_ALLOWED);
	}
}
