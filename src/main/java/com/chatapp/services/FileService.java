package com.chatapp.services;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;

import javax.servlet.http.HttpServletResponse;

public class FileService {
	private static FileService instance = null;

	public static Path rootLocation = Paths.get("archive");
	
	public static String rootURL = "http://localhost:8080/files/";

	private static final int DEFAULT_BUFFER_SIZE = 10240;

	private FileService() {

	}

	public static FileService getInstace() {
		if (instance == null) {
			instance = new FileService();
		}
		return instance;
	}

	public void handleStreamFileToClient(File file, String contentType, HttpServletResponse response) {
		if (contentType == null) {
			contentType = "application/octet-stream";
		}
		response.reset();
		response.setBufferSize(DEFAULT_BUFFER_SIZE);
		response.setContentType(contentType);
		response.setHeader("Content-Length", String.valueOf(file.length()));
		response.setHeader("Content-Disposition", "attachment; filename=\"" + file.getName() + "\"");
		BufferedInputStream input = null;
		BufferedOutputStream output = null;
		try {
			input = new BufferedInputStream(new FileInputStream(file), DEFAULT_BUFFER_SIZE);
			output = new BufferedOutputStream(response.getOutputStream(), DEFAULT_BUFFER_SIZE);

			byte[] buffer = new byte[DEFAULT_BUFFER_SIZE];
			int length;
			while ((length = input.read(buffer)) > 0) {
				output.write(buffer, 0, length);
			}
		} catch (IOException ex) {
			ex.printStackTrace();
		} finally {
			try {
				output.close();
				input.close();
			} catch (IOException ex) {
			}
		}
	}
}
