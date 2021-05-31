package com.chatapp.websockets;

import java.nio.ByteBuffer;
import java.util.LinkedList;
import java.util.Queue;

import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.PathParam;
import javax.websocket.server.ServerEndpoint;

import com.chatapp.models.MessageDecoder;
import com.chatapp.models.MessageEncoder;
import com.chatapp.models.dtos.FileDTO;
import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.services.ChatServiceAbstract;
import com.chatapp.services.MessageServiceInterface;
import com.chatapp.services.impl.ChatService;
import com.chatapp.services.impl.MessageService;

@ServerEndpoint(value = "/chat/{username}", encoders = MessageEncoder.class, decoders = MessageDecoder.class)
public class ChatWebsocket {

	private Session session;
	private String username;
	private Queue<FileDTO> fileDTOs = new LinkedList<>();

	private ChatServiceAbstract chatService = ChatService.getInstance();
	private MessageServiceInterface messageService = MessageService.getInstance();

	@OnOpen
	public void onOpen(@PathParam("username") String username, Session session) {
		if (chatService.register(this)) {
			this.session = session;
			this.username = username;
			String receiver = "all";
			MessageDTO msgResponse = new MessageDTO(this.username, "[P]open", "text", receiver, null);
			chatService.sendMessageToAllUsers(msgResponse);
		}
	}

	@OnError
	public void onError(Session session, Throwable throwable) {

	}

	@OnMessage
	public void onMessage(MessageDTO message, Session session) {
		chatService.sendMessageToOneUser(message, fileDTOs);
		messageService.saveMessage(message);
	}

	@OnMessage
	public void processUploading(ByteBuffer byteBuffer, boolean last, Session session) {
		System.err.println(byteBuffer.array().length);
		chatService.handleFileUpload(byteBuffer, last, fileDTOs);
	}

	@OnClose
	public void onClose(Session session) {
		if (chatService.close(this)) {
			String receiver = "all";
			MessageDTO msgResponse = new MessageDTO(this.username, "[P]close", "text", receiver, null);
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
