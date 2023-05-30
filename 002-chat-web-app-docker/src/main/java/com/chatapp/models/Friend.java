package com.chatapp.models;

public class Friend {
	private String sender;
	private String receiver;
	private String owner;
	private boolean status;

	public Friend() {

	}

	public Friend(String sender, String receiver, String owner, boolean status) {
		this.sender = sender;
		this.receiver = receiver;
		this.owner = owner;
		this.status = status;
	}

	public String getSender() {
		return sender;
	}

	public void setSender(String sender) {
		this.sender = sender;
	}

	public String getReceiver() {
		return receiver;
	}

	public void setReceiver(String receiver) {
		this.receiver = receiver;
	}

	public String getOwner() {
		return owner;
	}

	public void setOwner(String owner) {
		this.owner = owner;
	}

	public boolean isStatus() {
		return status;
	}

	public void setStatus(boolean status) {
		this.status = status;
	}
}
