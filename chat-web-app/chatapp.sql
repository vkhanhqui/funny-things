drop database chatapp;

create database chatapp COLLATE Latin1_General_100_CI_AI_KS_WS_SC;

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
	sender char(50),
	receiver char(50),
	owner char(50),
	status bit NOT NULL,
	foreign key (receiver) references users(username),
	foreign key (sender) references users(username),
	foreign key (owner) references users(username),
	primary key(sender, receiver, owner)
) ;

create table conversations(
	id int identity(1,1) primary key,
	name nvarchar(50) NOT NULL,
	avatar char(50) NOT NULL
) ;

create table conversations_users(
	conversations_id int,
	username char(50),
	is_admin bit NOT NULL,
	foreign key (conversations_id) references conversations(id),
	foreign key (username) references users(username),
	primary key(conversations_id, username)
) ;

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

--users
insert into users values('a1','a1',1,'a1.jpg');
insert into users values('a2','a2',1,'a2.jpg');
insert into users values('a3','a3',1,'a3.jpg');
insert into users values('a4','a4',1,'a4.jpg');

----friends
insert into friends values('a1','a2','a1',1);
insert into friends values('a2','a1','a1',1);

insert into friends values('a1','a3','a1',1);
insert into friends values('a3','a1','a1',1);

insert into friends values('a4','a2','a2',1);
insert into friends values('a2','a4','a2',1);

insert into friends values('a1','a4','a1',0);
insert into friends values('a4','a1','a1',0);

--message
--a1 vs a2
insert into messages(sender, receiver, message, message_type) 
	values('a1','a2','a1 chào a2','text');
--a2 vs a1
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 xin chào Võ Khánh Quí nà a1','text');
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 hello a1 - lan 2','text');
insert into messages(sender, receiver, message, message_type) 
	values('a2','a1','a2 hello a1 - lan 3','text');
insert into messages(sender, receiver, message, message_type) 
	values('a1','a2','a1 hello a2 - lan 2','text');
--a3 vs a1
insert into messages(sender, receiver, message, message_type) 
	values('a3','a1','a3 hello a1','text');
--a3 vs a2
insert into messages(sender, receiver, message, message_type) 
	values('a3','a2','a3 hello a2','text');

--conversation
insert into conversations(name, avatar) values('con heo'
	,concat('group-'
	,CAST(IDENT_CURRENT('conversations') as char(50))
	,'.png'
));

insert into conversations_users(conversations_id, username, is_admin) 
	values(IDENT_CURRENT('conversations'),'a1',1);
insert into conversations_users(conversations_id, username, is_admin) 
	values(IDENT_CURRENT('conversations'),'a2',0);
insert into conversations_users(conversations_id, username, is_admin) 
	values(IDENT_CURRENT('conversations'),'a3',0);

--conversation_user
--con heo group has a1, a2, a3 -> a1 is admin
insert into conversations_users(conversations_id, username, is_admin) 
	values(1,'a1',1);
insert into conversations_users(conversations_id, username, is_admin) 
	values(1,'a2',0);
insert into conversations_users(conversations_id, username, is_admin) 
	values(1,'a3',0);

--con ga group has a1, a2 -> a1 is admin
insert into conversations_users(conversations_id, username, is_admin) 
	values(2,'a1',1);
insert into conversations_users(conversations_id, username, is_admin) 
	values(2,'a2',0);

--con cho group has a1, a3 -> a1 is admin
insert into conversations_users(conversations_id, username, is_admin) 
	values(3,'a1',1);
insert into conversations_users(conversations_id, username, is_admin) 
	values(3,'a3',0);

--messages in conversation con heo
insert into messages(sender, message, message_type, conversations_id) 
	values('a1','a1 hello group con heo', 'text',1);
insert into messages(sender, message, message_type, conversations_id) 
	values('a2','a2 hello group con heo', 'text',1);
insert into messages(sender, message, message_type, conversations_id) 
	values('a2','a2 hello group con heo-lan 2', 'text',1);
insert into messages(sender, message, message_type, conversations_id) 
	values('a3','a3 hello group con heo', 'text',1);

--check
select * from users;

select * from friends;

select * from messages;

select * from conversations;

select * from conversations_users;

--load friends of user a1 (include requested friend)
select u2.username,u2.avatar,u2.gender 
from users u1 join friends f on u1.username = f.receiver 
join users u2 on u2.username = f.sender
where u1.username = 'a1';

--load messages of user a1 vs a2
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

--a1 find friends to add into a group
select u2.username,u2.avatar,u2.gender 
from users u1 join friends f on u1.username = f.receiver 
join users u2 on u2.username = f.sender
where u1.username = 'a1' 
and f.status = 1
and u2.username like '%a%'
and u2.username not in (
	select u.username
	from users u join conversations_users cu
		on u.username = cu.username
	join conversations c
		on c.id = cu.conversations_id
	where c.id = 1
);

--load groups which a1 is joined
select c.id, c.name
from conversations c join conversations_users cu
	on c.id = cu.conversations_id
where cu.username = 'a1';

--load users in con heo group
select u.username, u.avatar, u.gender
from users u join conversations_users cu
	on u.username = cu.username
join conversations c
	on c.id = cu.conversations_id
where c.id = 1

--load messages in con heo group 
select m.sender, u.avatar, m.message, m.receiver, m.message_type
from messages m join conversations c
on m.conversations_id = c.id
join users u on u.username = m.sender
where c.id = 1
order by created_at asc

--delete conversation by id
delete from conversations_users 
where conversations_id= 1;

delete from conversations 
where id = 1;

--delete user from conversation
delete from conversations_users 
where username = 'a2' 
and conversations_id = 1;

--find conversation by keyword
select c.id, c.name, c.avatar
from conversations c join conversations_users cu
on cu.conversations_id = c.id
where c.name like '%con heo%'
and cu.username = 'a1'

select * from conversations;