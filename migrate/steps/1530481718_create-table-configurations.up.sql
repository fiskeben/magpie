create table configurations (
    id serial primary key,
    service_name varchar(255) not null,
    version varchar(255),
    created_at timestamp,
    configuration jsonb
);
