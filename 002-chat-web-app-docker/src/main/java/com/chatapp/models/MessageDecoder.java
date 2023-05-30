package com.chatapp.models;

import java.io.IOException;

import javax.websocket.DecodeException;
import javax.websocket.Decoder;
import javax.websocket.EndpointConfig;

import com.chatapp.models.dtos.MessageDTO;
import com.fasterxml.jackson.databind.ObjectMapper;

public final class MessageDecoder implements Decoder.Text<MessageDTO> {

	private static ObjectMapper objectMapper = new ObjectMapper();

	@Override
	public void destroy() {
	}

	@Override
	public void init(final EndpointConfig arg0) {
	}

	@Override
	public MessageDTO decode(final String arg0) throws DecodeException {
		try {
			return objectMapper.readValue(arg0, MessageDTO.class);
		} catch (IOException e) {
			throw new DecodeException(arg0, "Unable to decode text to Message", e);
		}
	}

	@Override
	public boolean willDecode(final String arg0) {
		return arg0.contains("username") && arg0.contains("message");
	}
}
