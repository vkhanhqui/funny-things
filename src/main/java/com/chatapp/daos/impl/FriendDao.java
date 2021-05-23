package com.chatapp.daos.impl;

import java.util.List;

import com.chatapp.mappers.impl.FriendMapper;
import com.chatapp.models.Friend;

public class FriendDao extends GenericDao<Friend> {

	public void saveFriend(boolean isAccept, Friend friend) {
		String sender = friend.getSender();
		String receiver = friend.getReceiver();
		String owner = friend.getOwner();
		Boolean status = friend.isStatus();
		StringBuilder sql1 = new StringBuilder("insert into friends values(?,?,?,?)");
		StringBuilder sql2 = new StringBuilder("insert into friends values(?,?,?,?)");
		if (isAccept) {
			sql1 = new StringBuilder("update friends set status=? where sender = ? and receiver = ?");
			sql2 = new StringBuilder("update friends set status=? where sender = ? and receiver = ?");
			save(sql1.toString(), status, sender, receiver);
			save(sql2.toString(), status, receiver, sender);
		} else {
			save(sql1.toString(), sender, receiver, owner, status);
			save(sql2.toString(), receiver, sender, owner, status);
		}
	}

	public Friend findFriend(String sender, String receiver) {

		StringBuilder sql = new StringBuilder(
				"select sender,receiver, owner, status from friends where sender=? and receiver=?");

		List<Friend> friends = query(sql.toString(), new FriendMapper(), sender, receiver);
		return friends.isEmpty() ? null : friends.get(0);
	}
}
