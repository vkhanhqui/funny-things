create database chatapp;

use chatapp;

create login mylogin with password ='mylogin';

sp_changedbowner mylogin;

create table testing
(
	username char(20) primary key,
	password char(20) not null,
	gender bit not null,
	avatar char(25) not null
);

select * from testing;

truncate table testing;
drop table testing;

create table member
(
	code int primary key,
	username char(20) not null unique,
	name char(20) not null,
	e_password char(100) not null,
	k_password char(100) not null
);

create table friend
(
	friend_code int not null ,
	member_code int not null ,
	primary key(friend_code,member_code),
	foreign key(member_code) references member,
	foreign key(friend_code) references member
);
create table friend_request 
(
	sent_by int,
	sent_to int,
	primary key(sent_by,sent_to),
	foreign key(sent_by) references member,
	foreign key (sent_to) references member
);

create table message ( 
	code BIGINT primary key,
	message_date date not null,
	message_time time not null,
	from_code int not null, to_code int not null, 
	message varchar(500) not null, status char(1) not null, 
	foreign key(from_code) references member, 
	foreign key(to_code) references member 
); 

create table notification ( 
	code BIGINT primary key, 
	notification_date date not null, 
	notification_time time not null, 
	member_code int not null, 
	entity_code int not null, 
	notification_type int not null, 
	foreign key(member_code) references member 
);

create table closed_account
(
	code int not null unique,
	foreign key(code) references member 
);
