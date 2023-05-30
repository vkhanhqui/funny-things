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
		if (users != null) {
			if (conversation.getId() == null) {
				StringBuilder sqlCreateConversation = new StringBuilder();
				sqlCreateConversation.append("insert into conversations(name, avatar)");
				sqlCreateConversation.append(" values(?");
				sqlCreateConversation.append(",concat('group-'");
				sqlCreateConversation.append(",CAST(IDENT_CURRENT('conversations') as char(50))))");
				Long conversationId = save(sqlCreateConversation.toString(), conversation.getName());
				conversation.setId(conversationId);
			}
			users.forEach(user -> {
				StringBuilder sqlAddUserToConversation = new StringBuilder();
				sqlAddUserToConversation.append("insert into conversations_users(conversations_id, username,");
				sqlAddUserToConversation.append(" is_admin)");
				sqlAddUserToConversation.append(" values(?,?,?)");
				save(sqlAddUserToConversation.toString(), conversation.getId(), user.getUsername(), user.isAdmin());
			});
		} else {
			StringBuilder sql = new StringBuilder();
			sql.append("update conversations");
			if (!conversation.getAvatar().isEmpty()) {
				sql.append(" set name=?,avatar=?");
				sql.append(" where id=?");
				conversation.setAvatar(conversation.getAvatar().replaceAll(" ", ""));
				save(sql.toString(), conversation.getName(), conversation.getAvatar(), conversation.getId());
			} else {
				sql.append(" set name=?");
				sql.append(" where id=?");
				save(sql.toString(), conversation.getName(), conversation.getId());
			}
		}
	}

	@Override
	public List<Conversation> findAllConversationsByUsername(String username) {
		StringBuilder sql = new StringBuilder();
		sql.append("select c.id, c.name, c.avatar");
		sql.append(" from conversations c join conversations_users cu");
		sql.append(" on c.id = cu.conversations_id");
		sql.append(" where cu.username = ?");
		List<Conversation> conversations = query(sql.toString(), new ConversationMapper(), username);
		return conversations;
	}

	@Override
	public Conversation findConversationById(Long id) {
		StringBuilder sql = new StringBuilder();
		sql.append("select c.id, c.name, c.avatar");
		sql.append(" from conversations c");
		sql.append(" where c.id = ?");
		List<Conversation> conversations = query(sql.toString(), new ConversationMapper(), id);
		return conversations.isEmpty() ? null : conversations.get(0);
	}

	@Override
	public void deleteConversationById(Long id) {
		StringBuilder sql = new StringBuilder();
		sql.append("delete from conversations_users ");
		sql.append(" where conversations_id= ?;");

		sql.append("delete from messages ");
		sql.append(" where conversations_id= ?;");

		sql.append("delete from conversations");
		sql.append(" where id = ?;");
		save(sql.toString(), id, id, id);
	}

	@Override
	public void deleteUserFromConversation(Long conversationId, String username) {
		StringBuilder sql = new StringBuilder();
		sql.append("delete from conversations_users");
		sql.append(" where conversations_id = ?");
		sql.append(" and username= ?;");
		save(sql.toString(), conversationId, username);
	}

	@Override
	public List<Conversation> findConversationsOfUserByKeyword(String username, String keyword) {
		StringBuilder sql = new StringBuilder();
		sql.append("select c.id, c.name, c.avatar");
		sql.append(" from conversations c join conversations_users cu");
		sql.append(" on cu.conversations_id = c.id");
		sql.append(" where c.name like ?");
		sql.append(" and cu.username = ?");
		String param = "%" + keyword + "%";
		List<Conversation> conversations = query(sql.toString(), new ConversationMapper(), param, username);
		return conversations;
	}

}
