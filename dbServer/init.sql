set names utf8mb4;

grant all on cyberon_chatroom.* to 'asd' @'172.%' identified by '1234';
grant all on cyberon_chatroom.* to 'asd' @'localhost' identified by '1234';
flush privileges;
create database if not exists `cyberon_chatroom` character set 'utf8mb4';

use cyberon_chatroom;

create table if not exists member (
    id int primary key auto_increment,
    email varchar(100) unique key not null,
    password char(129) not null,
    nickname varchar(255) not null,
    created datetime default current_timestamp,
    updated datetime default current_timestamp
);

create table if not exists chatroom (
    id int primary key auto_increment,
    title varchar(100) unique key not null,
    member_id int,
    created datetime default current_timestamp,
    updated datetime default current_timestamp,
    foreign key(member_id) references member(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
