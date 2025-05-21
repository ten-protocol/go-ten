INSERT INTO results_type (uuid, environment, type, time, outcome)
VALUES (
    UUID(),
    'ten.sepolia',
    'test',
    UNIX_TIMESTAMP(),
    FALSE
);


select * from results_type where environment='ten.sepolia' and type='test'  order by time desc limit 2;

