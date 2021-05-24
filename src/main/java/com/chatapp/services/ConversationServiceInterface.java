package com.chatapp.services;

import java.util.List;

import com.chatapp.models.dtos.ConversationDTO;

public interface ConversationServiceInterface {
	public void saveConversation(ConversationDTO conversationDTO);

	public List<ConversationDTO> getAllConversationsByUsername(String username);
}
