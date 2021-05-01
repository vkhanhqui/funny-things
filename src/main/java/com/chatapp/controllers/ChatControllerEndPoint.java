package com.chatapp.controllers;

import java.util.Objects;

import javax.websocket.CloseReason.CloseCodes;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.PathParam;
import javax.websocket.server.ServerEndpoint;

import com.chatapp.models.Constants;
import com.chatapp.models.Message;
import com.chatapp.services.ChatService;
import com.chatapp.utils.MessageDecoder;
import com.chatapp.utils.MessageEncoder;

@ServerEndpoint(value = "/chat/{userId}", encoders = MessageEncoder.class, decoders = MessageDecoder.class)
public class ChatControllerEndPoint {

	@OnOpen
	public void onOpen(@PathParam(Constants.USERID_KEY) final String userId, final Session session) {
		if (Objects.isNull(userId) || userId.isEmpty()) {
			throw new RegistrationFailedException("userId is required");
		} else {
			if (ChatService.register(session)) {
				session.getUserProperties().put(Constants.USERID_KEY, userId);
				ChatService.onlineList.add(userId);
				System.out.printf("Session opened for %s\n", userId);
				String receiver = "all";
				ChatService.publish(new Message((String) session.getUserProperties().get(Constants.USERID_KEY),
						"[P]open", receiver, ChatService.onlineList), session);
			} else {
				throw new RegistrationFailedException("Unable to register, userId already exists, try another");
			}
		}
	}

	@OnError
	public void onError(final Session session, final Throwable throwable) {
		if (throwable instanceof RegistrationFailedException) {
			ChatService.close(session, CloseCodes.VIOLATED_POLICY, throwable.getMessage());
		}
	}

	@OnMessage
	public void onMessage(final Message message, final Session session) {
		ChatService.publish(message, session);
	}

	@OnClose
	public void onClose(final Session session) {
		if (ChatService.remove(session)) {
			System.out.printf("Session closed for %s\n", session.getUserProperties().get(Constants.USERID_KEY));

			String receiver = "all";
			ChatService.publish(new Message((String) session.getUserProperties().get(Constants.USERID_KEY),
					"[P]close", receiver, ChatService.onlineList), session);
		}
	}

	private static final class RegistrationFailedException extends RuntimeException {

		private static final long serialVersionUID = 1L;

		public RegistrationFailedException(final String message) {
			super(message);
		}
	}
}
