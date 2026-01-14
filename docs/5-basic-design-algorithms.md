# Overall design decisions

We can't use UUID or Snowflake solutions for ID/key generation because they need 16 byte
key with 12 character base64 string that it's not acceptable. For 1.2 billions link we can
use incremental ID in UINT32 datatype but high availablity requirements don't let us thinking about
databases with strong consistency and isolation. We choose BASE solutions to reduce failover downtime
to achive high availability. BASE database do not support autoincrement ID and we need a lock-free
solution for generating 32-64 bit keys to assignging in long urls.

# Shortener

Shortener service is an stateless application with n instance. In shortening link service will
use Cassandra LWT to counting a counter and generate unique keys per node.

## Key generation
System need a lock-free solution because BASE database do not support strong consistency.
Considered solution is keeping a table on Cassandra/ScyllaDB and use LWT for allocating ranges
to nodes. Each instance should have its own key range and work independently.

Following code is the base idea for generating keys.

```go
package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"math/bits"
	"os"

	"github.com/gocql/gocql"
)

/*
ScyllaDB INIT
CREATE KEYSPACE keygen WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'} ;
CREATE TABLE keygen.counters (key int PRIMARY KEY, value bigint);
*/
func main() {
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	))

	cluster := gocql.NewCluster(
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	)
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		slog.Error("failed to create scylladb session", "error", err)
		os.Exit(-1)
	}
	slog.Debug("scylladb connected successfully")

	for clusterID := range 20 {
		startCounter := clusterID * ((1 << 31) - 1)
		err = session.Query(
			`INSERT INTO keygen.counters (key, value) VALUES (?, ?) IF NOT EXISTS`,
			clusterID, startCounter).Exec()
		if err != nil {
			slog.Error("failed to init counter table", "error", err)
			os.Exit(1)
		}

		go func(clusterID int) {
			for range 10000 {
				keys, err := allocateRange(session, clusterID)
				if err != nil {
					slog.Error("failed to get keys", "error", err)
					continue
				}
				for _, k := range keys {
					fmt.Println(k)
				}
			}
		}(clusterID)
	}

	select {}
}

var KeyRange uint64 = 50

// CREATE TABLE keygen.counters (key int PRIMARY KEY, value bigint);
func allocateRange(session *gocql.Session, id int) ([]string, error) {
	var n uint64
	err := session.Query("SELECT value FROM keygen.counters WHERE key = ?", id).Scan(&n)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter value: %w", err)
	}

	applied, err := session.Query(
		`UPDATE keygen.counters SET value = ? WHERE key = ? IF value = ?`, n+KeyRange, id, n,
	).ScanCAS(&n)
	if err != nil {
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	if !applied {
		return nil, errors.New("do not applied")
	}

	r := []string{}
	for i := n; i < n+KeyRange; i++ {
		r = append(r, GenerateID(i))
	}
	return r, nil
}

func GenerateID(counter uint64) string {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, counter)
	nbits := (bits.Len64(counter) + 7) >> 3
	nbits = max(nbits, 4)
	return base64.RawURLEncoding.EncodeToString(key[:nbits])
}
```

## Data System
Shortener implelemt a port/adapter design pattern for accessing database, cache and telemetry.
We can change these layers in future.

## Database
Each link creation and lookup will request to BASE database.

## Cache
Service will cache all created links for one day and cache visited links for just an hour.

## Telemetry
Logs expose for all actions and handle later in observability systems.
BI need enumarate creation link, visited link telemetry.
Finally we need check how many times a link visited, calculate traffic peak hours
for autoscaling strategies and business solutions, performance metrics like
hardware usage and response time in various percentiles.
Keeping number of expired links is important for changing time of expiration links
and reporting to businesses how many users need its links event after expiration.

### GC pauses
Using key pool will due to GC pause and we can manage it will distributing service and minimazing meemory usage.

# 8 byte key
32 bit integer is sufficient for 6 character base64 encoded string. But we have a big limit in scaling key generator.
If we choose 32 bit key and limit all keys to 6 characters we should accept 3 key generator node limit. Why?
We should assign maximum capacity of of 5 years links (1.2 Bilion) for each key generator node.
3 node limit is not acceptable. System can utilize 64 bit counter and start with first 32 bit.
In the critical situations for backup and disaster recovery plan, nothing wrong in generating some 7-12 character keys!

We can utilize such a algorithm for generating keys:

```go
func GenerateID(counter uint64) string {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, counter)
	nbits := (bits.Len64(counter) + 7) >> 3
	return base64.RawURLEncoding.EncodeToString(key[:nbits])
}
```

With this algorithm I decide start from 6 character and if scale became large we can add more characters.
Also key expiration is the option.
