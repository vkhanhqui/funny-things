package com.chatapp.daos;

import java.util.List;

import com.chatapp.models.Conversation;
import com.chatapp.models.User;

public interface ConversationDaoInterface {
	void saveConversation(Conversation conversation, List<User> users);

	List<Conversation> findAllConversationsByUsername(String username);

	Conversation findConversationById(Long id);
	
	List<Conversation> findConversationsOfUserByKeyword(String username, String keyword);

	void deleteConversationById(Long id);

	void deleteUserFromConversation(Long conversationId, String username);
}
