package com.chatapp.websockets;

import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.PathParam;
import javax.websocket.server.ServerEndpoint;

import com.chatapp.models.Message;
import com.chatapp.services.ChatService;
import com.chatapp.utils.MessageDecoder;
import com.chatapp.utils.MessageEncoder;

@ServerEndpoint(value = "/chat/{username}", encoders = MessageEncoder.class, decoders = MessageDecoder.class)
public class ChatWebsocket {

	private Session session;

	private String username;

	private ChatService chatService = ChatService.getInstance();

	@OnOpen
	public void onOpen(@PathParam("username") String username, Session session) {
		if (chatService.register(this)) {
			this.session = session;
			this.username = username;
			String receiver = "all";
			Message msgResponse = new Message(this.username, "[P]open", receiver);
			chatService.sendMessageToAllUsers(msgResponse);
		}
	}

	@OnError
	public void onError(Session session, Throwable throwable) {

	}

	@OnMessage
	public void onMessage(Message message, Session session) {
		chatService.sendMessageToOneUser(message);
	}

	@OnClose
	public void onClose(Session session) {
		if (chatService.close(this)) {
			String receiver = "all";
			Message msgResponse = new Message(this.username, "[P]close", receiver);
			chatService.sendMessageToAllUsers(msgResponse);
		}
	}

	public Session getSession() {
		return session;
	}

	public void setSession(Session session) {
		this.session = session;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}
}
