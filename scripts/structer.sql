create table shortened_urls (
    id SERIAL primary key,
    url_base text unique not null,
    url_short text unique not null
);