package com.chatapp.services;

import java.util.List;

import javax.servlet.http.Part;

import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.models.dtos.UserDTO;

public interface ConversationServiceInterface {
	public void saveConversation(ConversationDTO conversationDTO);

	public List<ConversationDTO> getAllConversationsByUsername(String username);

	public List<UserDTO> getAllUsersByConversationId(Long id);

	public List<MessageDTO> getAllMessagesByConversationId(Long id);

	public void updateConversationById(Long id, String name, Part avatar);

	public ConversationDTO getConversationById(Long id);

	void deleteConversationById(Long id);

	void deleteUserFromConversation(Long conversationId, String username);

	public List<ConversationDTO> getConversationsOfUserByKeyword(String username, String keyword);
}
