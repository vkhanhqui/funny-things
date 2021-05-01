package com.chatapp.repositories;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class DBConnection {

	public static Connection connection;

	public static boolean Open() {
		try {
            Class.forName("com.microsoft.sqlserver.jdbc.SQLServerDriver");
			String url = "jdbc:sqlserver://KHANHQUI\\SQLEXPRESS:1433;databaseName=chatapp;user=mylogin;password=mylogin";
			connection = DriverManager.getConnection(url);
		} catch (SQLException | ClassNotFoundException ex) {
			return false;
		}
		return true;
	}

	public static boolean Close() {
		try {
			connection.close();
		} catch (SQLException ex) {
			return false;
		}
		return true;
	}

}
