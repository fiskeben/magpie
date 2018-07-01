create table blacklists (
    id serial primary key,
    service_name varchar(255),
    name varchar(255) not null
);
