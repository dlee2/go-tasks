DROP TABLE IF EXISTS task;
DROP TABLE IF EXISTS status;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS user;

CREATE TABLE task (
    id integer primary key autoincrement,
    title varchar(100),
    content text,
    task_status_id references status(id),
    created_date timestamp,
    last_modified_at timestamp,
    priority integer,
    cat_id references category(id),
    user_id references user(id),
    hide int);
CREATE TABLE status (
    id integer primary key autoincrement,
    status varchar(50) not null
);
CREATE TABLE files(
    name varchar(1000) not null,
    autoName varchar(255) not null,
    user_id references user(id),
    created_date timestamp
);
CREATE TABLE category(
    id integer primary key autoincrement,
    name varchar(1000) not null,
    user_id references user(id)
);

CREATE TABLE user (
    id integer primary key autoincrement,
    username varchar(100),
    password varchar(1000),
    email varchar(100)
);

insert into status(status) values('COMPLETE');
insert into status(status) values('PENDING');
insert into status(status) values('DELETED');
