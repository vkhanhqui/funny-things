package com.chatapp.services;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;

import javax.servlet.http.HttpServletResponse;

public abstract class FileServiceAbstract {

	public static Path rootLocation = Paths.get("archive");
	
	public static String rootURL = "http://localhost:8080/files/";

	protected static final int DEFAULT_BUFFER_SIZE = 10240;

	public abstract void handleStreamFileToClient(File file, String contentType, HttpServletResponse response);
}
