package com.chatapp.controllers;

import java.io.File;
import java.io.IOException;
import java.net.URLDecoder;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.chatapp.services.FileServiceAbstract;
import com.chatapp.services.impl.FileService;

@WebServlet("/files/*")
public class FileController extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private FileServiceAbstract fileService = FileService.getInstace();

	public FileController() {
		super();
	}

	protected void doGet(HttpServletRequest request, HttpServletResponse response)
			throws ServletException, IOException {
		String requestedFile = request.getPathInfo();
		if (requestedFile == null) {
			response.sendError(HttpServletResponse.SC_NOT_FOUND);
		} else {
			String filePath = FileServiceAbstract.rootLocation.toString();
			File file = new File(filePath, URLDecoder.decode(requestedFile, "UTF-8"));
			if (!file.exists()) {
				response.sendError(HttpServletResponse.SC_NOT_FOUND);
			} else {
				String contentType = getServletContext().getMimeType(file.getName());
				fileService.handleStreamFileToClient(file, contentType, response);
			}
		}
	}
}