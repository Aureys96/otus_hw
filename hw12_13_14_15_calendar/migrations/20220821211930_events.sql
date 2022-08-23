CREATE SEQUENCE events_id_sequence
    START 1
    INCREMENT 1;

create table events (
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('events_id_sequence'),
    title varchar(256) not null,
    start_at timestamp,
    end_at timestamp,
    user_id int,
    created_at timestamp default now()
);
