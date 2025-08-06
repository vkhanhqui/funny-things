create database chatapp COLLATE Latin1_General_100_CI_AI_KS_WS_SC;
go

use chatapp;
go

CREATE TABLE users (
	username char(50) primary key,
	password char(50) NOT NULL,
	gender bit NOT NULL,
	avatar char(50) NOT NULL
) ;
go

CREATE TABLE friends (
	sender char(50),
	receiver char(50),
	owner char(50),
	status bit NOT NULL,
	foreign key (receiver) references users(username),
	foreign key (sender) references users(username),
	foreign key (owner) references users(username),
	primary key(sender, receiver, owner)
) ;
go

create table conversations(
	id int identity(1,1) primary key,
	name nvarchar(50) NOT NULL,
	avatar char(50) NOT NULL
) ;
go

create table conversations_users(
	conversations_id int,
	username char(50),
	is_admin bit NOT NULL,
	foreign key (conversations_id) references conversations(id),
	foreign key (username) references users(username),
	primary key(conversations_id, username)
) ;
go

CREATE TABLE messages (
	id int identity(1,1) primary key,
	sender char(50) NOT NULL,
	receiver char(50),
	message nvarchar(max) NOT NULL,
	message_type char(100) NOT NULL,
	created_at datetime default current_timestamp,
	conversations_id int,
	foreign key (sender) references users(username),
	foreign key (receiver) references users(username),
	foreign key (conversations_id) references conversations(id)
) ;
go