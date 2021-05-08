package com.chatapp.models;

import java.io.IOException;

import javax.websocket.DecodeException;
import javax.websocket.Decoder;
import javax.websocket.EndpointConfig;

import com.fasterxml.jackson.databind.ObjectMapper;

public final class MessageDecoder implements Decoder.Text<Message> {

	@Override
	public void destroy() {
	}

	@Override
	public void init(final EndpointConfig arg0) {
	}

	@Override
	public Message decode(final String arg0) throws DecodeException {
		try {
			return new ObjectMapper().readValue(arg0, Message.class);
		} catch (IOException e) {
			throw new DecodeException(arg0, "Unable to decode text to Message", e);
		}
	}

	@Override
	public boolean willDecode(final String arg0) {
		return arg0.contains("username") && arg0.contains("message");
	}
}
