create schema if not exists messenger;

create table if not exists messenger.messages (
    id 				uuid 	        primary key not null,
    sender_id 		bigint 	        not null,
    chat_id      	varchar(255) 	not null,
    text 			text 	        not null
);