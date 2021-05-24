package com.chatapp.daos.impl;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.ResourceBundle;

import com.chatapp.daos.GenericDaoInterface;
import com.chatapp.mappers.RowMapperInterface;

public class GenericDao<T> implements GenericDaoInterface<T> {

	private ResourceBundle resourceBundle = ResourceBundle.getBundle("database");

	public Connection getConnection() {
		try {
			String driverName = resourceBundle.getString("driverName");
			String server = resourceBundle.getString("server");
			String databaseName = resourceBundle.getString("databaseName");
			String user = resourceBundle.getString("user");
			String password = resourceBundle.getString("password");
			StringBuilder url = new StringBuilder();
			Class.forName(driverName);
			url.append(server);
			url.append(";databaseName=" + databaseName);
			url.append(";user=" + user);
			url.append(";password=" + password);
//			url.append(";integratedSecurity=true");
			return DriverManager.getConnection(url.toString());
		} catch (ClassNotFoundException | SQLException ex) {
			return null;
		}
	}

	@Override
	public List<T> query(String sql, RowMapperInterface<T> rowMapper, Object... parameters) {
		List<T> results = new ArrayList<>();
		Connection connection = null;
		PreparedStatement statement = null;
		ResultSet resultSet = null;
		try {
			connection = getConnection();
			statement = connection.prepareStatement(sql);
			setParameter(statement, parameters);
			resultSet = statement.executeQuery();
			while (resultSet.next()) {
				results.add(rowMapper.mapRow(resultSet));
			}
			return results;
		} catch (SQLException ex) {
			return new ArrayList<>();
		} finally {
			try {
				if (connection != null) {
					connection.close();
				}
				if (statement != null) {
					statement.close();
				}
				if (resultSet != null) {
					resultSet.close();
				}
			} catch (SQLException ex) {
				return new ArrayList<>();
			}
		}
	}

	private void setParameter(PreparedStatement prepareStatement, Object... parameters) {
		try {
			for (int i = 0; i < parameters.length; i++) {
				Object parameter = parameters[i];
				int index = i + 1;
				if (parameter instanceof Long) {
					prepareStatement.setLong(index, (Long) parameter);
				} else if (parameter instanceof String) {
					prepareStatement.setString(index, (String) parameter);
				} else if (parameter instanceof Integer) {
					prepareStatement.setInt(index, (Integer) parameter);
				} else if (parameter instanceof Timestamp) {
					prepareStatement.setTimestamp(index, (Timestamp) parameter);
				} else if (parameter instanceof Boolean) {
					prepareStatement.setBoolean(index, (Boolean) parameter);
				}
			}
		} catch (SQLException ex) {
			ex.printStackTrace();
		}
	}

	@Override
	public Long save(String sql, Object... parameters) {
		Connection connection = null;
		PreparedStatement preparedStatement = null;
		ResultSet resultSet = null;
		try {
			Long id = null;
			connection = getConnection();
			connection.setAutoCommit(false);
			preparedStatement = connection.prepareStatement(sql, Statement.RETURN_GENERATED_KEYS);
			setParameter(preparedStatement, parameters);
			preparedStatement.executeUpdate();
			resultSet = preparedStatement.getGeneratedKeys();
			if (resultSet.next()) {
				id = resultSet.getLong(1);
			}
			connection.commit();
			return id;
		} catch (SQLException ex) {
			if (connection != null) {
				try {
					connection.rollback();
				} catch (SQLException ex1) {
					ex1.printStackTrace();
				}
			}
		} finally {
			try {
				if (connection != null) {
					connection.close();
				}
				if (preparedStatement != null) {
					preparedStatement.close();
				}
			} catch (SQLException ex2) {
				ex2.printStackTrace();
			}
		}
		return null;
	}

	@Override
	public int count(String sql, Object... parameters) {
		Connection connection = null;
		PreparedStatement statement = null;
		ResultSet resultSet = null;
		try {
			int count = 0;
			connection = getConnection();
			statement = connection.prepareStatement(sql);
			setParameter(statement, parameters);
			resultSet = statement.executeQuery();
			while (resultSet.next()) {
				count = resultSet.getInt(1);
			}
			return count;
		} catch (SQLException ex) {
			return 0;
		} finally {
			try {
				if (connection != null) {
					connection.close();
				}
				if (statement != null) {
					statement.close();
				}
				if (resultSet != null) {
					resultSet.close();
				}
			} catch (SQLException ex) {
				return 0;
			}
		}
	}

}
