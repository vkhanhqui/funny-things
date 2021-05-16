drop database chatapp;

create database chatapp;

use chatapp;

create login mylogin with password ='mylogin';

sp_changedbowner mylogin;

CREATE TABLE users (
	username char(50) primary key,
	password char(50) NOT NULL,
	gender bit NOT NULL,
	avatar char(50) NOT NULL
) ;

CREATE TABLE friends (
	sender char(50) NOT NULL,
	foreign key (sender) references users(username),
	receiver char(50) NOT NULL,
	foreign key (receiver) references users(username),
	primary key(sender, receiver),
	status bit not null
) ;

CREATE TABLE messages (
	id int identity(1,1) primary key,
	sender char(50) NOT NULL,
	foreign key (sender) references users(username),
	receiver char(50) NOT NULL,
	foreign key (receiver) references users(username),
	message text NOT NULL,
	message_type char(20) not null,
	created_at datetime default current_timestamp
) ;

insert into users values('a1','a1',1,'a1.jpg');
insert into users values('a2','a2',1,'a2.jpg');
insert into users values('a3','a3',1,'a3.jpg');
insert into users values('a4','a4',1,'a4.jpg');

insert into friends values('a1','a2',1);
insert into friends values('a2','a1',1);

insert into friends values('a1','a3',1);
insert into friends values('a3','a1',1);

insert into friends values('a4','a2',1);
insert into friends values('a2','a4',1);

insert into friends values('a1','a4',0);
insert into friends values('a4','a1',0);

insert into messages(sender, receiver, message, message_type) 
	values('a1','a2','a1 hello a2','text');
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 hello a1','text');
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 hello a1 - lan 2','text');
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 hello a1 - lan 3','text');
insert into messages(sender, receiver, message, message_type) 
	values('a1','a2','a1 hello a2 - lan 2','text');

insert into messages(sender, receiver, message, message_type) 
	values('a3','a1','a3 hello a1','text');
insert into messages(sender, receiver, message, message_type) 
	values('a3','a2','a3 hello a2','text');

select * from users;

select * from friends;

--load friends
select u2.username,u2.avatar,u2.gender 
from users u1 join friends f on u1.username = f.receiver 
join users u2 on u2.username = f.sender
where u1.username LIKE 'a1' and f.status = 1;

--load messages
select m1.sender, m1.message, m1.message_type, m1.receiver
from messages m1 inner join(
		select id
		from messages
		where sender = 'a2'
		or receiver ='a2'
		) m2 on m1.id = m2.id
where m1.sender = 'a1'
or m1.receiver = 'a1'
order by created_at asc

truncate table friends;
truncate table messages;
truncate table users;

drop table messages;