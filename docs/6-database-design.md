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

# Observability system

I need complete observability. We have many chooices but I select OTel collector with
VictoriaMetrics for metrics, Loki for logs and Jeager for traces.

- OTel: Collector
- VictoriaMetrisc: Metrics
- Loki: Logs
- Jeager: Traces
- Grafana: Visualizing