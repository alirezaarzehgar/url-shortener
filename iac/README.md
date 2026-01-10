# Setup quick ScyllaDB cluster

set `fs.aio-max-nr` to a big number to run ScyllaDB successfully

```bash
sudo sysctl fs.aio-max-nr=675900
```

Then start ScyllaDB cluster

```bash
docker compose -f iac/docker-compose.yml up -d
docker exec -it scylla1 nodetool status
docker exec -it scylla1 cqlsh
```
