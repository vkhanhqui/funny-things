package com.chatapp.daos;

import java.util.List;

import com.chatapp.mappers.RowMapperInterface;

public interface GenericDaoInterface<T> {
	List<T> query(String sql, RowMapperInterface<T> rowMapper, Object... parameters);

	Long save(String sql, Object... parameters);

	int count(String sql, Object... parameters);
}
