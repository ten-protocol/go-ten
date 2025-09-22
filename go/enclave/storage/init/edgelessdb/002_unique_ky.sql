SET SESSION rocksdb_max_row_locks = 10000000;
SET SESSION net_read_timeout = 600;
SET SESSION net_write_timeout = 600;

DELETE s1 FROM statedb32 s1 INNER JOIN statedb32 s2 ON s1.ky = s2.ky AND s1.id > s2.id;
ALTER TABLE statedb32 ADD CONSTRAINT unique_ky_statedb32 UNIQUE (ky);

DELETE FROM statedb64  WHERE id NOT IN (SELECT MIN(id) FROM statedb64 GROUP BY ky);
ALTER TABLE statedb64 ADD CONSTRAINT unique_ky_statedb64 UNIQUE (ky);