package com.chatapp.services;

import com.chatapp.models.dtos.ConversationDTO;

public interface ConversationServiceInterface {
	public void saveConversation(ConversationDTO conversationDTO);
}
