drop database if exists bbs;
create database bbs;
use bbs;

drop table if exists users;
drop table if exists threads;
drop table if exists posts;
drop table if exists sessions;
drop table if exists thread_class;

create table users (
  id         int primary key auto_increment,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null
) engine=innodb default charset=utf8;

create table thread_class (
  id         int primary key auto_increment,
  name       varchar(64) not null unique
) engine=innodb default charset=utf8;
insert into thread_class values(1, "其他");

create table threads (
  id         int primary key auto_increment,
  uuid       varchar(64) not null unique,
  topic      text,
  body       test,
  user_id    int not null,
  class_id   int default 0,
  created_at timestamp not null,
  foreign key(user_id) references users(id),
  foreign key(class_id) references thread_class(id)
) engine=innodb default charset=utf8;
                               
create table posts (           
  id         int primary key auto_increment,
  uuid       varchar(64) not null unique,
  body       text,
  user_id    int not null,
  thread_id  int not null,
  created_at timestamp not null,
  foreign key(user_id) references users(id),
  foreign key(thread_id) references threads(id)
) engine=innodb default charset=utf8;
               
create table sessions (
  id         int primary key auto_increment,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    int not null,
  created_at timestamp not null,
  foreign key(user_id) references users(id)
) engine=innodb default charset=utf8;
