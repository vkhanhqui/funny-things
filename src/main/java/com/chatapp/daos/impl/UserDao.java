package com.chatapp.daos.impl;

import java.util.List;
import java.util.stream.Collectors;

import com.chatapp.daos.UserDaoInterface;
import com.chatapp.mappers.impl.UserMapper;
import com.chatapp.models.User;

public class UserDao extends GenericDao<User> implements UserDaoInterface {

	private static UserDao instance = null;

	private UserDao() {

	}

	public synchronized static UserDao getInstace() {
		if (instance == null) {
			instance = new UserDao();
		}
		return instance;
	}

	@Override
	public User findByUserNameAndPassword(String userName, String password) {
		StringBuilder sql = new StringBuilder("select username, gender, avatar");
		sql.append(" from users where username=? and password=?");
		List<User> users = query(sql.toString(), new UserMapper(), userName, password);
		return users.isEmpty() ? null : users.get(0);
	}

	@Override
	public List<User> findFriends(String userName) {
		StringBuilder sql = new StringBuilder("select distinct u2.username, u2.avatar, u2.gender");
		sql.append(" from users u1 join friends f on u1.username = f.receiver");
		sql.append(" join users u2 on u2.username = f.sender");
		sql.append(" where u1.username LIKE ?");
		String param = "%" + userName + "%";
		List<User> users = query(sql.toString(), new UserMapper(), param);
		return users.stream().filter(u -> !u.getUsername().equals(userName)).collect(Collectors.toList());
	}

	@Override
	public void saveUser(Boolean isRegister, User user) {
		String username = user.getUsername();
		String password = user.getPassword();
		Boolean gender = user.isGender();
		String avatar = user.getAvatar();
		StringBuilder sql = new StringBuilder("insert into users values(?,?,?,?)");
		if (isRegister) {
			save(sql.toString(), username, password, gender, avatar);
		} else {
			sql = new StringBuilder("update users set password=?, gender=?, avatar=? where username=?");
			save(sql.toString(), password, gender, avatar, username);
		}
	}

	@Override
	public List<User> findFriendsByKeyWord(String userName, String keyWord) {
		StringBuilder sql = new StringBuilder("select u2.username, u2.avatar, u2.gender");
		sql.append(" from users u2 where username != ? and username like ?");
		String param = "%" + keyWord + "%";
		List<User> users = query(sql.toString(), new UserMapper(), userName, param);
		return users;

	}

	@Override
	public List<User> findUsersByConversationId(Long id) {
		StringBuilder sql = new StringBuilder();
		sql.append("select u.username, u.avatar, u.gender, cu.is_admin");
		sql.append(" from users u join conversations_users cu");
		sql.append(" on u.username = cu.username");
		sql.append(" join conversations c");
		sql.append(" on c.id = cu.conversations_id");
		sql.append(" where c.id = ?");
		List<User> users = query(sql.toString(), new UserMapper(), id);
		return users;
	}

	@Override
	public List<User> findFriendsNotInConversation(String userName, String keyword, Long conversationId) {
		StringBuilder sql = new StringBuilder();
		sql.append("select u2.username,u2.avatar,u2.gender");
		sql.append(" from users u1 join friends f on u1.username = f.receiver ");
		sql.append(" join users u2 on u2.username = f.sender");
		sql.append(" where u1.username = ?");
		sql.append(" and f.status = 1");
		sql.append(" and u2.username like ?");
		sql.append(" and u2.username not in (");
		sql.append(" select u.username");
		sql.append(" from users u join conversations_users cu");
		sql.append(" on u.username = cu.username");
		sql.append(" join conversations c");
		sql.append(" on c.id = cu.conversations_id");
		sql.append(" where c.id = ?)");
		String param = "%" + keyword + "%";
		List<User> users = query(sql.toString(), new UserMapper(), userName, param, conversationId);
		return users;
	}

}
