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
| key | bigint | Generated numeral key | PRIMARY KEY |
| original_url | Origianl URl for redirection | NOT NULL |
| created_at | TIMESTAMP | Creation timestamp | DEFAULT CURRENT_TIMESTAMP |

Key generation table schema

| Column | Type | Description | Constraint |
| --- | --- | --- | --- |
| cluster_id | int | Worker process id assigned in condig | PRIMARY KEY |
| counter | bigint | Key counter for generating base64 key  | DEFUALT 0 |


## ScyllaDB migration
```sql
-- migrate
CREATE KEYSPACE urlshortener WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};
CREATE TABLE urlshortener.keypool (key int PRIMARY KEY, counter bigint);
CREATE TABLE urlshortener.urls (key text PRIMARY KEY, original_url text, created_at timestamp);
-- Seed
INSERT INTO urlshortener.keypool (key, counter) VALUES (1, 0);
```

# Observability system

I need complete observability. We have many chooices but I select OTel collector with
VictoriaMetrics for metrics, Loki for logs and Jeager for traces.

- OTel: Collector
- VictoriaMetrisc: Metrics
- Loki: Logs
- Jeager: Traces
- Grafana: Visualizing