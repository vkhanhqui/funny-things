package com.chatapp.services.impl;

import java.util.List;
import java.util.stream.Collectors;

import com.chatapp.daos.ConversationDaoInterface;
import com.chatapp.daos.impl.ConversationDao;
import com.chatapp.models.Conversation;
import com.chatapp.models.User;
import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.models.dtos.UserDTO;
import com.chatapp.services.ConversationServiceInterface;

public class ConversationService implements ConversationServiceInterface {

	private ConversationDaoInterface conversationDaoInterface = ConversationDao.getInstance();

	private static ConversationService instance = null;

	private ConversationService() {

	}

	public synchronized static ConversationService getInstance() {
		if (instance == null) {
			instance = new ConversationService();
		}
		return instance;
	}

	private User convertToUserEntity(UserDTO userDTO) {
		User user = new User();
		user.setUsername(userDTO.getUsername());
		user.setAdmin(userDTO.isAdmin());
		return user;
	}

	@Override
	public void saveConversation(ConversationDTO conversationDTO) {
		Conversation conversation = new Conversation();
		conversation.setName(conversationDTO.getName());
		List<User> users = conversationDTO.getUsers().stream().map(userDTO -> convertToUserEntity(userDTO))
				.collect(Collectors.toList());
		conversationDaoInterface.saveConversation(conversation, users);
	}

}
