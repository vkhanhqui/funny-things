package com.chatapp.daos.impl;

import java.util.List;

import com.chatapp.daos.ConversationDaoInterface;
import com.chatapp.mappers.impl.ConversationMapper;
import com.chatapp.models.Conversation;
import com.chatapp.models.User;

public class ConversationDao extends GenericDao<Conversation> implements ConversationDaoInterface {

	private static ConversationDao instance = null;

	private ConversationDao() {

	}

	public synchronized static ConversationDao getInstance() {
		if (instance == null) {
			instance = new ConversationDao();
		}
		return instance;
	}

	@Override
	public void saveConversation(Conversation conversation, List<User> users) {
		StringBuilder sqlCreateConversation = new StringBuilder();
		sqlCreateConversation.append("insert into conversations(name)");
		sqlCreateConversation.append(" values(?)");
		Long conversationId = save(sqlCreateConversation.toString(), conversation.getName());
		conversation.setId(conversationId);

		users.forEach(user -> {
			StringBuilder sqlAddUserToConversation = new StringBuilder();
			sqlAddUserToConversation.append("insert into conversations_users(conversations_id, username,");
			sqlAddUserToConversation.append(" is_admin)");
			sqlAddUserToConversation.append(" values(?,?,?)");
			save(sqlAddUserToConversation.toString(), conversationId, user.getUsername(), user.isAdmin());
		});
	}

	@Override
	public List<Conversation> findAllConversationsByUsername(String username) {
		StringBuilder sql = new StringBuilder();
		sql.append("select c.id, c.name");
		sql.append(" from conversations c join conversations_users cu");
		sql.append(" on c.id = cu.conversations_id");
		sql.append(" where cu.username = ?");
		List<Conversation> conversations = query(sql.toString(), new ConversationMapper(), username);
		return conversations;
	}

}
