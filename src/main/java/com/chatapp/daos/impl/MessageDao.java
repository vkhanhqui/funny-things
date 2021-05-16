package com.chatapp.daos.impl;

import java.util.List;

import com.chatapp.daos.MessageDaoInterface;
import com.chatapp.mappers.impl.MessageMapper;
import com.chatapp.models.Message;

public class MessageDao extends GenericDao<Message> implements MessageDaoInterface {

	private static MessageDao instance = null;

	private MessageDao() {

	}

	public synchronized static MessageDao getInstance() {
		if (instance == null) {
			instance = new MessageDao();
		}
		return instance;
	}

	@Override
	public List<Message> findAllMessagesBySenderAndReceiver(String sender, String receiver) {
		StringBuilder sql = new StringBuilder("select m1.sender, m1.message, m1.message_type, m1.receiver");
		sql.append(" from messages m1 inner join(");
		sql.append("select id from messages");
		sql.append(" where sender = ? or receiver = ? )");
		sql.append(" m2 on m1.id = m2.id");
		sql.append(" where m1.sender = ? ");
		sql.append(" or m1.receiver = ? ");
		sql.append(" order by created_at asc");
		List<Message> listMessages = query(sql.toString(), new MessageMapper(), receiver, receiver, sender, sender);
		return listMessages;
	}

	@Override
	public void saveMessage(Message message) {
		StringBuilder sql = new StringBuilder("insert into messages(sender, receiver, message, message_type)");
		sql.append(" values(?,?,?,?)");
		String sender = message.getUsername();
		String msg = message.getMessage();
		String type = message.getType();
		String receiver = message.getReceiver();
		save(sql.toString(), sender, receiver, msg, type);
	}

}
