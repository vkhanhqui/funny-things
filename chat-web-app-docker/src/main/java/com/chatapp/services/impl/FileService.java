package com.chatapp.services.impl;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;

import javax.servlet.http.HttpServletResponse;

import com.chatapp.services.FileServiceAbstract;

public class FileService extends FileServiceAbstract{
	private static FileService instance = null;

	private FileService() {

	}

	public static FileService getInstace() {
		if (instance == null) {
			instance = new FileService();
		}
		return instance;
	}

	@Override
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
