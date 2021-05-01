package com.chatapp.services;

import java.io.IOException;
import java.util.HashSet;
import java.util.Objects;
import java.util.Set;
import java.util.concurrent.CopyOnWriteArraySet;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import javax.websocket.CloseReason;
import javax.websocket.CloseReason.CloseCodes;

import com.chatapp.models.Constants;
import com.chatapp.models.Message;

import javax.websocket.EncodeException;
import javax.websocket.Session;

public class ChatService {

	private static final Lock LOCK = new ReentrantLock();
	private static final Set<Session> SESSIONS = new CopyOnWriteArraySet<>();
	public static Set<String> onlineList = new HashSet<String>();

	private ChatService() {
		throw new IllegalStateException(Constants.INSTANTIATION_NOT_ALLOWED);
	}

	public static void publish(Message message, final Session origin) {
		assert !Objects.isNull(message) && !Objects.isNull(origin);
		if (!message.getReceiver().equals("all")) {
			SESSIONS.stream().filter(session -> session.getUserProperties().containsValue(message.getReceiver()))
					.forEach(session -> {
						try {
							session.getBasicRemote().sendObject(message);
						} catch (IOException | EncodeException e) {
							e.printStackTrace();
						}
					});
		} else {
			SESSIONS.stream().forEach(session -> {
				try {
					session.getBasicRemote().sendObject(message);
				} catch (IOException | EncodeException e) {
					e.printStackTrace();
				}
			});
		}
	}

	public static boolean register(final Session session) {
		assert !Objects.isNull(session);

		boolean result = false;
		try {
			LOCK.lock();

			result = !SESSIONS.contains(session) && !SESSIONS.stream()
					.filter(elem -> ((String) elem.getUserProperties().get(Constants.USERID_KEY))
							.equals((String) session.getUserProperties().get(Constants.USERID_KEY)))
					.findFirst().isPresent() && SESSIONS.add(session);
		} finally {
			LOCK.unlock();
		}

		return result;
	}

	public static void close(final Session session, final CloseCodes closeCode, final String message) {
		assert !Objects.isNull(session) && !Objects.isNull(closeCode);

		try {
			session.close(new CloseReason(closeCode, message));
		} catch (IOException e) {
			throw new RuntimeException("Unable to close session", e);
		}
	}

	public static boolean remove(final Session session) {
		assert !Objects.isNull(session);
		onlineList.remove(session.getUserProperties().get(Constants.USERID_KEY).toString());
		return SESSIONS.remove(session);
	}
}
