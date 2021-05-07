package com.chatapp.websockets;

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

@ServerEndpoint(value = "/chat/{username}", encoders = MessageEncoder.class, decoders = MessageDecoder.class)
public class ChatWebsocket {
	private ChatService chatService = ChatService.getInstance();

	@OnOpen
	public void onOpen(@PathParam(Constants.USERNAME_KEY) final String username, final Session session) {
		if (Objects.isNull(username) || username.isEmpty()) {
			throw new RegistrationFailedException("userId is required");
		} else {
			if (chatService.register(session)) {
				session.getUserProperties().put(Constants.USERNAME_KEY, username);
				ChatService.onlineList.add(username);
				System.out.printf("Session opened for %s\n", username);
				String receiverResponse = "all";
				String usernameResponse = (String) session.getUserProperties().get(Constants.USERNAME_KEY);
				Message msgResponse = new Message(usernameResponse, "[P]open", receiverResponse,
						ChatService.onlineList);
				chatService.publish(msgResponse, session);
			} else {
				throw new RegistrationFailedException("Unable to register, userId already exists, try another");
			}
		}
	}

	@OnError
	public void onError(final Session session, final Throwable throwable) {
		if (throwable instanceof RegistrationFailedException) {
			chatService.close(session, CloseCodes.VIOLATED_POLICY, throwable.getMessage());
		}
	}

	@OnMessage
	public void onMessage(final Message message, final Session session) {
		chatService.publish(message, session);
	}

	@OnClose
	public void onClose(final Session session) {
		if (chatService.remove(session)) {
			System.out.printf("Session closed for %s\n", session.getUserProperties().get(Constants.USERNAME_KEY));

			String receiverResponse = "all";
			String usernameResponse = (String) session.getUserProperties().get(Constants.USERNAME_KEY);
			Message msgResponse = new Message(usernameResponse, "[P]close", receiverResponse, ChatService.onlineList);
			chatService.publish(msgResponse, session);
		}
	}

	private static final class RegistrationFailedException extends RuntimeException {

		private static final long serialVersionUID = 1L;

		public RegistrationFailedException(final String message) {
			super(message);
		}
	}
}
