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
import com.chatapp.services.ChatSessionManager;
import com.chatapp.utils.MessageDecoder;
import com.chatapp.utils.MessageEncoder;

@ServerEndpoint(value = "/chat-controller/{email}", encoders = MessageEncoder.class, decoders = MessageDecoder.class)
public final class ChatControllerEndPoint {

	@OnOpen
	public void onOpen(@PathParam(Constants.EMAIL_KEY) final String email, final Session session) {
		if (Objects.isNull(email) || email.isEmpty()) {
			throw new RegistrationFailedException("Email is required");
		} else {
			session.getUserProperties().put(Constants.EMAIL_KEY, email);
			if (ChatSessionManager.register(session)) {
				System.out.printf("Session opened for %s\n", email);

				ChatSessionManager
						.publish(new Message((String) session.getUserProperties().get(Constants.EMAIL_KEY),
								"***joined the chat***"), session);
			} else {
				throw new RegistrationFailedException("Unable to register, email already exists, try another");
			}
		}
	}

	@OnError
	public void onError(final Session session, final Throwable throwable) {
		if (throwable instanceof RegistrationFailedException) {
			ChatSessionManager.close(session, CloseCodes.VIOLATED_POLICY, throwable.getMessage());
		}
	}

	@OnMessage
	public void onMessage(final Message message, final Session session) {
		ChatSessionManager.publish(message, session);
	}

	@OnClose
	public void onClose(final Session session) {
		if (ChatSessionManager.remove(session)) {
			System.out.printf("Session closed for %s\n", session.getUserProperties().get(Constants.EMAIL_KEY));

			ChatSessionManager.publish(new Message((String) session.getUserProperties().get(Constants.EMAIL_KEY),
					"***left the chat***"), session);
		}
	}

	private static final class RegistrationFailedException extends RuntimeException {

		private static final long serialVersionUID = 1L;

		public RegistrationFailedException(final String message) {
			super(message);
		}
	}
}
