package com.chatapp.services;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.CopyOnWriteArraySet;

import javax.websocket.EncodeException;

import com.chatapp.models.Message;
import com.chatapp.websockets.ChatWebsocket;

public class ChatService {

	private static ChatService chatService = null;
	private static final Set<ChatWebsocket> chatWebsockets = new CopyOnWriteArraySet<>();
	private Queue<FileOutputStream> fileOutputStreams = new LinkedList<>();
	private Queue<String> receivers = new LinkedList<>();

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
		if (!message.getType().equals("text")) {
			String fileName = message.getMessage();
			String destFile = RegisterService.rootLocation.toString() + "/" + message.getUsername() + "/" + fileName;
			File uploadedFile = new File(destFile);
			if (!uploadedFile.exists()) {
				try {
					FileOutputStream fileOutputStream = new FileOutputStream(uploadedFile);
					fileOutputStreams.add(fileOutputStream);
					receivers.add(message.getReceiver());
				} catch (FileNotFoundException ex) {
					ex.printStackTrace();
				}
			}
		} else {
			chatWebsockets.stream().filter(chatWebsocket -> chatWebsocket.getUsername().equals(message.getReceiver()))
					.forEach(chatWebsocket -> {
						try {
							chatWebsocket.getSession().getBasicRemote().sendObject(message);
						} catch (IOException | EncodeException e) {
							e.printStackTrace();
						}
					});
		}
	}

	public void handleFileUpload(String username, ByteBuffer byteBuffer, boolean last) {
		try {
			if (!last) {
				while (byteBuffer.hasRemaining()) {
					fileOutputStreams.peek().write(byteBuffer.get());
				}
			} else {
				fileOutputStreams.peek().flush();
				fileOutputStreams.peek().close();
				fileOutputStreams.remove();
				String message = "alo alo";
				String type = "text";
				String receiver = receivers.peek();
				Message messageResponse = new Message(username, message, type, receiver);
				sendMessageToOneUser(messageResponse);
				receivers.remove();
			}

		} catch (IOException ex) {
			System.out.println(ex.getMessage());
			ex.printStackTrace();
		}
	}

	private Set<String> getUsernames() {
		Set<String> usernames = new HashSet<String>();
		chatWebsockets.forEach(chatWebsocket -> {
			usernames.add(chatWebsocket.getUsername());
		});
		return usernames;
	}
}
