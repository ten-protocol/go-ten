SELECT 
    t.time, t.gasused,
    (t.gasprice * t.gasused * r.rate / 1e18) AS ten_sepolia
FROM 
    tx_costs t
JOIN 
    (SELECT rate FROM rates WHERE currency = 'USD' ORDER BY time DESC LIMIT 1) r
ON 1=1 -- cross join
WHERE 
    t.name = 'Game.deploy' 
    AND t.environment = 'ten.sepolia'
ORDER by time DESC;