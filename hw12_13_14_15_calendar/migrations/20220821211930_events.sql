create table events (
    id serial primary key,
    title varchar(256) not null,
    start_at timestamp,
    end_at timestamp,
    user_id int,
    created_at timestamp default now()
);
