create database messenger1;

\connect messenger1

create extension if not exists "uuid-ossp";

create table if not exists db_id (
	exists bool primary key default true,
	id uuid not null,
    constraint singlerow check (exists)
);

insert into db_id(id) values (uuid_generate_v4());


create database messenger2;

\connect messenger2

create extension if not exists "uuid-ossp";

create table if not exists db_id (
	exists bool primary key default true,
	id uuid not null,
    constraint singlerow check (exists)
);

insert into db_id(id) values (uuid_generate_v4());