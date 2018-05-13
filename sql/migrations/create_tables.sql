create table roasters (
  id         serial constraint roaster_id_key primary key,
  name       varchar(255),
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  website    varchar(2083)
);
create unique index roasters_id_idx
  on roasters (id);
create index roasters_name_idx
  on roasters (name);
create index roasters_created_at_idx
  on roasters (created_at);
create index roasters_updated_at_idx
  on roasters (updated_at);

create table coffees (
  id         serial constraint coffees_id_key primary key,
  name       varchar(1023) not null,
  created_at timestamp     not null default current_timestamp,
  updated_at timestamp     not null default current_timestamp,
  roaster_id integer references roasters (id)
);
create unique index coffees_id_idx
  on coffees (id);
create index coffees_name_idx
  on coffees (name);
create index coffees_created_at_idx
  on coffees (created_at);
create index coffees_updated_at_idx
  on coffees (updated_at);
create index coffees_roaster_id_idx
  on coffees (roaster_id);

create table reviews (
  id          serial constraint reviews_id_key primary key,
  reviewed_at timestamp not null,
  created_at  timestamp not null default current_timestamp,
  updated_at  timestamp not null default current_timestamp,
  review_link varchar(2083),
  roaster_id  integer references roasters (id),
  coffee_id   integer references coffees (id),
  score       integer
);
create unique index reviews_id_idx
  on reviews (id);
create index reviewed_at_idx
  on reviews (reviewed_at);
create index reviews_created_at_idx
  on reviews (created_at);
create index reviews_updated_at_idx
  on reviews (updated_at);
create index reviews_roaster_id_idx
  on reviews (roaster_id);
create index reviews_coffee_id_idx
  on reviews (coffee_id);
create index score
  on reviews (score);
