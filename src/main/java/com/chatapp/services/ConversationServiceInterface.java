package com.chatapp.services;

import java.util.List;

import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.models.dtos.UserDTO;

public interface ConversationServiceInterface {
	public void saveConversation(ConversationDTO conversationDTO);

	public List<ConversationDTO> getAllConversationsByUsername(String username);
	
	public List<UserDTO> getAllUsersByConversationId(Long id);
	
	public List<MessageDTO> getAllMessagesByConversationId(Long id);
}
