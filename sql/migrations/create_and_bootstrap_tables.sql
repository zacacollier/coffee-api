/* create `roasters` table */
drop table if exists roasters cascade;
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


/* create `coffees` table */
drop table if exists coffees cascade;
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

/* create `reviews` table */
drop table if exists reviews cascade;
create table reviews (
  id          serial constraint reviews_id_key primary key,
  reviewed_at timestamp not null,
  created_at  timestamp not null default current_timestamp,
  updated_at  timestamp not null default current_timestamp,
  review_link varchar(2083),
  roaster_id  integer references roasters (id),
  coffee_id   integer references coffees (id),
  score       integer            default null
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

/* Provision temp table */
drop table if exists temp_csv_imports;
CREATE TABLE temp_csv_imports (
  id           serial constraint temp_csv_imports_id_key primary key,
  rating       varchar(4),
  coffee       varchar(1023),
  roaster      varchar(255),
  review_link  varchar(2083),
  roaster_link varchar(2083),
  excerpt      text,
  review_date  varchar(31)
);
/* need to run this from psql for now */
copy temp_csv_imports (
  rating,
  coffee,
  roaster,
  review_link,
  roaster_link,
  excerpt,
  review_date
)
FROM '/home/zac/go/src/github.com/zacacollier/CoffeeAPI/coffee_reviews.csv' (
FORMAT CSV );
update temp_csv_imports
set rating = null
where rating = 'NR';

/* bootstrap `roasters` */
insert into roasters (name, website)
  select distinct
    roaster      as name,
    roaster_link as website
  from temp_csv_imports;

/* bootstrap `coffees` */
insert into coffees (name, roaster_id)
  select distinct
    tci.coffee as name,
    r.id       as roaster_id
  from temp_csv_imports tci
    left join roasters r on r.name = tci.roaster;

/* bootstrap `reviews` */
insert into reviews (reviewed_at, review_link, roaster_id, coffee_id, score)
  select distinct
    tci.review_date :: timestamp as reviewed_at,
    tci.review_link,
    r.id                         as roaster_id,
    c.id                         as coffee_id,
    (
      case when tci.rating = 'NR'
        then null
      else tci.rating :: INTEGER
      end
    )                            as score
  from
    temp_csv_imports tci
    left join coffees c on c.name = tci.coffee
    left join roasters r on r.id = c.roaster_id;
