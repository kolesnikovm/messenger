-- +goose Up
create table if not exists messenger.messages (
    id 				uuid 	        primary key not null,
    sender_id 		bigint 	        not null,
    chat_id      	varchar(255) 	not null,
    text 			text 	        not null
);

-- +goose Down
drop table if exists messenger.messages;
