-- +goose Up
create schema if not exists messenger;

-- +goose Down
drop schema if exists messenger;
