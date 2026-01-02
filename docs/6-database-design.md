# Shortener service
## Database Choice
Shortener service need a persistent key values store to handle 1:100 write:read traffic. We don't need joins
and we consider BASE databases. We have some options with key value data model option and highly avalablity
with eventual consistency and quorom for tolerating node outages without failover. Leaderless replication is
the considered solution.

Options:
- Apache Cassandra
- ScyllaDB
- Amazon DynamoDB
- Azure CosmosDB
- Riak KV

In this options DynamoDB and CosmosDB are very niche. Riak KV is good but have poor community right now
but from data model perspective it's perfect. Also ScyllaDB are prefered over Cassandra because less
hardware usage, lack of GC pauses, less operational overhead, better documentation and community.

Among this options I choose ScyllaDB. But everything is changable in the future. Next option is Riak KV.

## Schema
Shortener service just need an schema with key/value structure.

| Column | Type | Description | Constraint |
| --- | --- | --- | --- |
| key | UINT64 | Generated numeral key | PRIMARY KEY |
| original_url | Origianl URl for redirection | NOT NULL |
| created_at | TIMESTAMP | Creation timestamp | DEFAULT CURRENT_TIMESTAMP |
| ttl | TTL | Expiration time | Less than 3 months |

# Key Generation Service
## Database Choice
Unique key generation need strong consistency. PostgreSQL is overkill for an small KGS. System need a database
that replicate easily, very lightweight with small failover time and overhead. We also need persistent key value store
but with strong consistency.

Options:
- etcd
- Redis with RDB
- ACID databases (overkill and high overhead)
- Cockroachdb
- RocksDB

We have two architecture:
- Managing database as server (etcd, redis db, pg/mysql): High operational overhead
- Developing KGS using embedded KV databases indpendently and loadbalancing traffic to them
  and config specific range for all of them. It will offer a super simple independent KGS systems
  with high fault tolerance.

Winner is embedded databases for KGS and scale it linearly. Then we can hide this KGSs over LB.

Options:
- **BadgerDB**
- **BoltDB** (now bbolt)
- **LevelDB** (via goleveldb)
- **Pebble**
- **RocksDB** (via gorocksdb)
- **SQLite** (as embedded KV using key-value tables)
- **TiKV** (in embedded/TiDB Local mode)
- **Bitcask** (via implementations like go-bitcask)
- **LMDB** (via bolt-mdb or go-lmdb)
- **NutsDB**
- **MosDB** (Moss store)
- **VictoriaMetrics' lib/tsdb** (for time-series-like keys)
- **BBolt** (fork of BoltDB)
- **Pogreb**

The choice is **Badger**.

## Schema

| Column | Type | Description | Constraints |
| --- | --- | --- | --- |
| key_counter | UINT64 | Counter for assigning and counting key ranges for KGS | NOT NULL |

# Observability system

I need complete observability. We have many chooices but I select OTel collector with
VictoriaMetrics for metrics, Loki for logs and Jeager for traces.

- OTel: Collector
- VictoriaMetrisc: Metrics
- Loki: Logs
- Jeager: Traces
- Grafana: Visualizing