package com.chatapp.models;

import com.fasterxml.jackson.databind.ObjectMapper;

public final class Constants {

	public static final String INSTANTIATION_NOT_ALLOWED = "Instantiation not allowed";
	public static final String EMAIL_KEY = "email";
	public static final String MESSAGE_KEY = "message";
	public static final ObjectMapper MAPPER = new ObjectMapper();

	private Constants() {
		throw new IllegalStateException(INSTANTIATION_NOT_ALLOWED);
	}
}
