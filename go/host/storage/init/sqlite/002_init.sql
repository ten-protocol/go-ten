create table if not exists block_count
(
id          int  NOT NULL PRIMARY KEY,
total       int  NOT NULL
);

insert into block_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;