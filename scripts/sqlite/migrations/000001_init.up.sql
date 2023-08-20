create table trigger
(
    id                   TEXT,
    enabled              INT,
    created              integer default current_timestamp,
    name                 TEXT,
    description          TEXT,
    gpiopin              integer,
    seconds_to_retrigger integer
);

create table webhook
(
    id         TEXT,
    trigger_id TEXT
        constraint webhook_trigger_id_fk
            references trigger (id),
    URL        integer,
    headers    BLOB,
    body       TEXT
);

