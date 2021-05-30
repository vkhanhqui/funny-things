package com.chatapp.services.impl;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.List;
import java.util.stream.Collectors;

import javax.servlet.http.Part;

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
import com.chatapp.services.FileServiceAbstract;

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
		conversationDTO.setAvatar(conversation.getAvatar().trim());
		return conversationDTO;
	}

	private Conversation convertToConversation(ConversationDTO conversationDTO) {
		Conversation conversation = new Conversation();
		conversation.setId(conversationDTO.getId());
		conversation.setName(conversationDTO.getName());
		if (conversationDTO.getAvatar() != null && !conversationDTO.getAvatar().isEmpty()) {
			conversation.setAvatar(conversationDTO.getAvatar().trim());
		}
		return conversation;
	}

	private MessageDTO convertToMessageDTO(Message messageEntity) {
		String username = messageEntity.getUsername();
		String type = messageEntity.getType();
		String message = messageEntity.getMessage();
		if (!type.equals("text")) {
			message = FileServiceAbstract.toTagHtml(type, username, message);
		}
		String receiver = messageEntity.getReceiver();
		Long groupId = messageEntity.getGroupId();
		MessageDTO messageDTO = new MessageDTO(username, message, type, receiver, groupId);
		messageDTO.setAvatar(messageEntity.getAvatar());
		return messageDTO;
	}

	@Override
	public void saveConversation(ConversationDTO conversationDTO) {
		Conversation conversation = convertToConversation(conversationDTO);
		List<User> users = conversationDTO.getUsers().stream().map(userDTO -> convertToUserEntity(userDTO))
				.collect(Collectors.toList());
		conversationDaoInterface.saveConversation(conversation, users);
		conversationDTO.setId(conversation.getId());

		String dirName = "group-" + conversationDTO.getId();
		File privateDir = new File(FileServiceAbstract.rootLocation.toString() + "/" + dirName);
		privateDir.mkdir();
		String fileName = dirName + ".png";
		File newFile = new File(privateDir.toString() + "/" + fileName);
		try {
			File defaultAvatar = new File(FileServiceAbstract.rootLocation.toString() + "/default/group.png");
			Files.copy(defaultAvatar.toPath(), newFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
			conversation.setAvatar(fileName);
			conversationDaoInterface.saveConversation(conversation, null);
			conversationDTO.setAvatar(fileName);
		} catch (IOException ex) {
		}
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

	@Override
	public void updateConversationById(Long id, String name, Part avatar) {
		try {
			String fileName = "";
			String origin = avatar.getSubmittedFileName();
			if (!origin.isEmpty()) {
				String dirName = "group-" + id;
				File privateDir = new File(FileServiceAbstract.rootLocation.toString() + "/" + dirName);
				String tail = origin.substring(origin.lastIndexOf("."), origin.length());
				fileName = dirName + tail;
				System.err.println("file: " + fileName);
				avatar.write(privateDir.getAbsolutePath() + File.separator + fileName);
			}
			Conversation conversation = new Conversation(id, name, fileName);
			conversationDaoInterface.saveConversation(conversation, null);
		} catch (IOException ex) {
		}
	}

	@Override
	public ConversationDTO getConversationById(Long id) {
		Conversation conversation = conversationDaoInterface.findConversationById(id);
		return convertToConversationDTO(conversation);
	}

	@Override
	public void deleteConversationById(Long id) {
		conversationDaoInterface.deleteConversationById(id);

	}

	@Override
	public void deleteUserFromConversation(Long conversationId, String username) {
		conversationDaoInterface.deleteUserFromConversation(conversationId, username);
	}

	@Override
	public List<ConversationDTO> getConversationsOfUserByKeyword(String username, String keyword) {
		List<ConversationDTO> conversationDTOs = conversationDaoInterface
				.findConversationsOfUserByKeyword(username, keyword).stream()
				.map(conversation -> convertToConversationDTO(conversation)).collect(Collectors.toList());
		return conversationDTOs;
	}
}
