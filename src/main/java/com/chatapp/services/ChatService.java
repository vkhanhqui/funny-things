package com.chatapp.services;

import java.io.IOException;
import java.util.Objects;
import java.util.Set;
import java.util.concurrent.CopyOnWriteArraySet;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import javax.websocket.CloseReason;
import javax.websocket.CloseReason.CloseCodes;
import javax.websocket.EncodeException;
import javax.websocket.Session;

import com.chatapp.models.Constants;
import com.chatapp.models.Message;

public class ChatService {

	private static ChatService chatService = null;
	private static Lock lock = new ReentrantLock();
	private static Set<Session> sessions = new CopyOnWriteArraySet<Session>();

	public static Set<String> onlineList = new CopyOnWriteArraySet<String>();

	private ChatService() {
	}

	public synchronized static ChatService getInstance() {
		if (chatService == null) {
			chatService = new ChatService();
		}
		return chatService;
	}

	public void publish(Message message, final Session origin) {
		assert !Objects.isNull(message) && !Objects.isNull(origin);
		if (!message.getReceiver().equals("all")) {
			sessions.stream().filter(session -> session.getUserProperties().containsValue(message.getReceiver()))
					.forEach(session -> {
						try {
							session.getBasicRemote().sendObject(message);
						} catch (IOException | EncodeException e) {
							e.printStackTrace();
						}
					});
		} else {
			sessions.stream().forEach(session -> {
				try {
					session.getBasicRemote().sendObject(message);
				} catch (IOException | EncodeException e) {
					e.printStackTrace();
				}
			});
		}
	}

	public boolean register(final Session session) {
		assert !Objects.isNull(session);

		boolean result = false;
		try {
			lock.lock();

			result = !sessions.contains(session) && !sessions.stream()
					.filter(elem -> ((String) elem.getUserProperties().get(Constants.USERNAME_KEY))
							.equals((String) session.getUserProperties().get(Constants.USERNAME_KEY)))
					.findFirst().isPresent() && sessions.add(session);
		} finally {
			lock.unlock();
		}

		return result;
	}

	public void close(final Session session, final CloseCodes closeCode, final String message) {
		assert !Objects.isNull(session) && !Objects.isNull(closeCode);

		try {
			session.close(new CloseReason(closeCode, message));
		} catch (IOException e) {
			throw new RuntimeException("Unable to close session", e);
		}
	}

	public boolean remove(final Session session) {
		assert !Objects.isNull(session);
		onlineList.remove(session.getUserProperties().get(Constants.USERNAME_KEY).toString());
		return sessions.remove(session);
	}
}
