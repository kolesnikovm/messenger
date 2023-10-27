create schema if not exists messenger;

create table if not exists messenger.messages (
    id 				uuid 	primary key not null,
    sender_id 		bigint 	not null,
    recipient_id 	bigint 	not null,
    text 			text 	not null
);