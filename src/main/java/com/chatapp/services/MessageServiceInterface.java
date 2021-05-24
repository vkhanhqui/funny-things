package com.chatapp.services;

import java.util.List;

import com.chatapp.models.dtos.MessageDTO;

public interface MessageServiceInterface {
	public List<MessageDTO> getAllMessagesBySenderAndReceiver(String sender, String receiver);

	public void saveMessage(MessageDTO messageDTO);
}
