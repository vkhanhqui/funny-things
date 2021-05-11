package com.chatapp.services;

import java.nio.ByteBuffer;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.CopyOnWriteArraySet;

import com.chatapp.models.FileDTO;
import com.chatapp.models.Message;
import com.chatapp.websockets.ChatWebsocket;

public abstract class ChatServiceAbstract {

	protected static final Set<ChatWebsocket> chatWebsockets = new CopyOnWriteArraySet<>();
	
	public abstract boolean register(ChatWebsocket chatWebsocket);

	public abstract boolean close(ChatWebsocket chatWebsocket);

	public abstract void sendMessageToAllUsers(Message message);

	public abstract void sendMessageToOneUser(Message message, Queue<FileDTO> fileDTOs);

	public abstract void handleFileUpload(ByteBuffer byteBuffer, boolean last, Queue<FileDTO> fileDTOs);

	protected abstract Set<String> getUsernames();
}
