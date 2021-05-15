drop database chatapp;

create database chatapp;

use chatapp;

create login mylogin with password ='mylogin';

sp_changedbowner mylogin;

CREATE TABLE users (
  id int identity(1,1) primary key,
  username char(50) NOT NULL,
  password char(50) NOT NULL,
  first_name varchar(50) NOT NULL,
  last_name varchar(50) NOT NULL,
  gender bit NOT NULL,
  avatar char(50) NOT NULL
) ;

CREATE TABLE friends (
  id_sender int NOT NULL,
  foreign key (id_sender) references users(id),
  id_receiver int NOT NULL,
  foreign key (id_receiver) references users(id),
  primary key(id_sender, id_receiver),
  status bit not null
) ;

CREATE TABLE messages (
  id int identity(1,1) primary key,
  created_at datetime NOT NULL,
  message text NOT NULL,
  id_sender int NOT NULL,
  foreign key (id_sender) references users(id),
  id_receiver int NOT NULL,
  foreign key (id_receiver) references users(id)
) ;

truncate table friends;
truncate table users;

select * from users;
select * from friends;

insert into users(username, password, first_name, last_name, gender, avatar) 
	values('a1','a1','a1','a1',1,'a1.jpg');
insert into users(username, password, first_name, last_name, gender, avatar) 
	values('a2','a2','a2','a2',1,'a2.jpg');
insert into users(username, password, first_name, last_name, gender, avatar) 
	values('a3','a3','a3','a3',1,'a3.jpg');
insert into users(username, password, first_name, last_name, gender, avatar) 
	values('a4','a4','a4','a4',1,'a4.jpg');

insert into friends values(1,2,1);
insert into friends values(2,1,1);

insert into friends values(1,3,1);
insert into friends values(3,1,1);

insert into friends values(4,2,1);
insert into friends values(2,4,1);


select u1.username u1,u1.avatar u2_avt, f.status, u2.username u2,u2.avatar u2_avt
from users u1 join friends f on u1.id = f.id_receiver
join users u2 on u2.id = f.id_sender
where u1.username LIKE 'a2' or  u2.username LIKE 'a2';

select u2.username,u2.avatar,u2.gender 
from users u1 join friends f on u1.id = f.id_receiver 
join users u2 on u2.id = f.id_sender
where u1.username LIKE 'a1'