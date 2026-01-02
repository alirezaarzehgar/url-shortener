# Overall design decisions

We can't use UUID or Snowflake solutions for ID/key generation because they need 16 byte
key with 12 character base64 string that it's not acceptable. For 1.2 billions link we can
use incremental ID in UINT32 datatype but high availablity requirements don't let us thinking about
databases with strong consistency and isolation. We choose BASE solutions to reduce failover downtime
to achive high availability. In BASE databases we can't utilize autoincremention and we need a KGS
(Key Generation Service). Our solution to generating keys is Hi/Lo Allocator and this algorithm
needs autoincremention. For this we can need consistency and we can accept failover down time for
small KGS database.  Also we can up n cluster of KGS databases and use service discover to route
service to that clusters. Naturally each KGS cluster will allocate specific range of keys and
choosing UINT32 is a constraint for 4 cluster system in the end. We can choose BIGINT to break
constraints and design for super big number 18446744073 possible clusters! in this case we
have no problem with DDOS on key reservation and linear scaling is possible!

# Shortener

Shortener service is an stateless application with n instance. In shortening link service will
request to KGS and get a short key. Service will save long link with generated key as primary key
for future lookups. KGS guaranteed keys are unique and collision is impossible.

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
Interconnection with KGS is important.

## Bottleneck
### Shortener and KGS interconnection
Interconnection between shortener service and KGS will be bottleneck. Shortener can send request for
bulk key generation and keep a pool of keys each time. This solution will speedup system for link creation
but performance issues will be emerge when a new link request touch empty pool and response time
will be inconsistent. For solving this problem service can send async request when key pool reached to
an specific number of it's capacity and with this solution response times will be consistent.

### GC pauses
Using key pool will due to GC pause and we can manage it will distributing service and minimazing meemory usage.


# Key Generation Service
KGS accept request for new keys and response new unique keys. It should be response bulk keys
for minimizing requests between services. KGS use a cluster of base databases and configured for
a range of keys. Each KGS cluster work on 2B range of keys and it should be configurable.

## Data System
KGS need a persistent key values store with strong consistency and linearizability for couting keys
and guarantee uniqueness. It will cause downtime in failover and can be managed by launching multiple
instance of KGS system with different range of keys.
Hi/Lo allocator will miss some keys in the instance outage but it's acceptable because we have big capacity.

## Algorithm
Hi/Lo Allocator works simple by a counter. Each instance request for end number of counter, reserve it with
defined range and increment shared counter. For example instance start working, add RANGE to COUNT and
keep that range on memory. When each internal service requests for new key, it will count internal counter
atomicly. Automatic pool management is a good solution also for KGS to preventing unpredictable and inconsistent
response times for key reservation.

## Cache
In KGS each request to getting keys are unique but keeping keys should be manage in cache.
Saving in local cache is good but it will causes to missing more keys in the case of instance outage.
We can utilize distributed cache for tolerating key loss.

But for simplecity we don't use distributed cache and consider it overkill. Using local cache for reserving future
key range is efficient and siple enough.

We save two ranges for each instance. When first range filled we offer keys of second range and request to database
asyncly. External service will not sense of inconsistent response with this strategy.

## Telemetry
KGS need expose response time, hardware usage, managing key pool, database quey results (success or failure rates)
and external service connections.

## Bottlenecks
### Database connections
KGS can keep a range of key ranges to avoid response time inconsistency. This solution is shared with shortener service.
We have two solutions. KGS generate keys and keep in the distributed cache then other services can give their keys from
distributed cache service but it will couple tightly in specific distributed cache.
Each service can keep data in separate or at least theorically separate distributed caches and manage data.

KGS can keep key ranges on cache, manage pool concurrently. Key pool management is a shared and common pattern/module
with KGS and other services that need this service and will reduce requests between KGS and database and internal services
and KGS.
