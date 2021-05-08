package com.chatapp.services;

import java.io.IOException;
import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.CopyOnWriteArraySet;

import javax.websocket.EncodeException;

import com.chatapp.models.Message;
import com.chatapp.websockets.ChatWebsocket;

public class ChatService {

	private static ChatService chatService = null;
	private static final Set<ChatWebsocket> chatWebsockets = new CopyOnWriteArraySet<>();

	private ChatService() {
	}

	public synchronized static ChatService getInstance() {
		if (chatService == null) {
			chatService = new ChatService();
		}
		return chatService;
	}

	public boolean register(ChatWebsocket chatWebsocket) {
		return chatWebsockets.add(chatWebsocket);
	}

	public boolean close(ChatWebsocket chatWebsocket) {
		return chatWebsockets.remove(chatWebsocket);
	}

	public void sendMessageToAllUsers(Message message) {
		message.setOnlineList(getUsernames());
		chatWebsockets.stream().forEach(chatWebsocket -> {
			try {
				chatWebsocket.getSession().getBasicRemote().sendObject(message);
			} catch (IOException | EncodeException e) {
				e.printStackTrace();
			}
		});
	}

	public void sendMessageToOneUser(Message message) {
		chatWebsockets.stream().filter(chatWebsocket -> chatWebsocket.getUsername().equals(message.getReceiver()))
				.forEach(chatWebsocket -> {
					try {
						chatWebsocket.getSession().getBasicRemote().sendObject(message);
					} catch (IOException | EncodeException e) {
						e.printStackTrace();
					}
				});
	}

	private Set<String> getUsernames() {
		Set<String> usernames = new HashSet<String>();
		chatWebsockets.forEach(chatWebsocket -> {
			usernames.add(chatWebsocket.getUsername());
		});
		return usernames;
	}
}
