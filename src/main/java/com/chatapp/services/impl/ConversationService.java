package com.chatapp.services.impl;

import java.util.List;
import java.util.stream.Collectors;

import com.chatapp.daos.ConversationDaoInterface;
import com.chatapp.daos.MessageDaoInterface;
import com.chatapp.daos.UserDaoInterface;
import com.chatapp.daos.impl.ConversationDao;
import com.chatapp.daos.impl.MessageDao;
import com.chatapp.daos.impl.UserDao;
import com.chatapp.models.Conversation;
import com.chatapp.models.Message;
import com.chatapp.models.User;
import com.chatapp.models.dtos.ConversationDTO;
import com.chatapp.models.dtos.MessageDTO;
import com.chatapp.models.dtos.UserDTO;
import com.chatapp.services.ConversationServiceInterface;

public class ConversationService implements ConversationServiceInterface {

	private ConversationDaoInterface conversationDaoInterface = ConversationDao.getInstance();

	private UserDaoInterface userDaoInterface = UserDao.getInstace();

	private MessageDaoInterface messageDaoInterface = MessageDao.getInstance();

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

	private UserDTO convertToUserDTO(User user) {
		UserDTO userDTO = new UserDTO();
		userDTO.setUsername(user.getUsername());
		userDTO.setAvatar(user.getAvatar());
		userDTO.setAdmin(user.isAdmin());
		return userDTO;
	}

	private ConversationDTO convertToConversationDTO(Conversation conversation) {
		ConversationDTO conversationDTO = new ConversationDTO();
		conversationDTO.setId(conversation.getId());
		conversationDTO.setName(conversation.getName());
		return conversationDTO;
	}

	private Conversation convertToConversation(ConversationDTO conversationDTO) {
		Conversation conversation = new Conversation();
		conversation.setId(conversationDTO.getId());
		conversation.setName(conversationDTO.getName());
		return conversation;
	}

	private MessageDTO convertToMessageDTO(Message message) {
		MessageDTO messageDTO = new MessageDTO();
		messageDTO.setUsername(message.getUsername());
		messageDTO.setMessage(message.getMessage());
		messageDTO.setType(message.getType());
		return messageDTO;
	}

	@Override
	public void saveConversation(ConversationDTO conversationDTO) {
		Conversation conversation = convertToConversation(conversationDTO);
		List<User> users = conversationDTO.getUsers().stream().map(userDTO -> convertToUserEntity(userDTO))
				.collect(Collectors.toList());
		conversationDaoInterface.saveConversation(conversation, users);
		conversationDTO.setId(conversation.getId());
	}

	@Override
	public List<ConversationDTO> getAllConversationsByUsername(String username) {
		List<Conversation> conversations = conversationDaoInterface.findAllConversationsByUsername(username);
		List<ConversationDTO> conversationDTOs = conversations.stream()
				.map(conversation -> convertToConversationDTO(conversation)).collect(Collectors.toList());
		return conversationDTOs;
	}

	@Override
	public List<UserDTO> getAllUsersByConversationId(Long id) {
		List<UserDTO> userDTOs = userDaoInterface.findUsersByConversationId(id).stream()
				.map(user -> convertToUserDTO(user)).collect(Collectors.toList());
		return userDTOs;
	}

	@Override
	public List<MessageDTO> getAllMessagesByConversationId(Long id) {
		List<MessageDTO> messageDTOs = messageDaoInterface.findAllMessagesByConvesationId(id).stream()
				.map(message -> convertToMessageDTO(message)).collect(Collectors.toList());
		return messageDTOs;
	}

}
