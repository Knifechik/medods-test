-- up

CREATE TABLE tokens_table
(
    id                 uuid        not null default gen_random_uuid(),
    user_id            uuid        not null,
    hash_refresh_token text        not null,
    ip                 text        not null,
    created_at         timestamp not null default now(),
    updated_at         timestamp not null default now(),
    expires_at         timestamp   not null,

    primary key (id),
    unique (hash_refresh_token),
    unique (user_id)
);


-- down

DROP TABLE tokens_table;