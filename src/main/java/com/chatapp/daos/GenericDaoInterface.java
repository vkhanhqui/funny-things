package com.chatapp.daos;

import java.util.List;

import com.chatapp.mappers.RowMapperInterface;

public interface GenericDaoInterface<T> {
	List<T> query(String sql, RowMapperInterface<T> rowMapper, Object... parameters);
	void update (String sql, Object... parameters);
	void insert (String sql, Object... parameters);
	int count(String sql, Object... parameters);
}
