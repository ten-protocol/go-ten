### Querying the Postgres DB started by the local testnet

The launcher starts a Postgres container called `pg-ten` on the `node_network` Docker network.

- Connect with psql inside the running container:

```bash
docker exec -it pg-ten psql -U postgres -h localhost -d postgres
# list databases (non-templates)
\l
# connect to a node DB (created automatically, e.g. host_<node_id_lowercased>)
\c host_...
# run your SQL query
SELECT * FROM host_rollup WHERE hash = blah;
```

