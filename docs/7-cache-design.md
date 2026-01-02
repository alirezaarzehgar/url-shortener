# Cache
System needs cache for 4 situations:

## Shortener
- Distributed cache for created urls for 1 day
- Distributed cache for visited urls for 1 hour
- Local cache for pool of generated key ranges

## Distributed Cache Schema
Only shortener needs distributed cache and needed schema is like this:

| Key Pattern | Key Type | Value Type | Description | TTL |
| --- | --- | --- | --- | --- |
| link:{key} | STRING | STRING | Original URL cache | 1 Day |
| visit:{key} | STRING | STRING | Visit URL cache after first day | 1 Hour |

### Technology choice

We have two options, Redis and Memcached. Memcached is simpler, lighter and managing
its cluster is easier that Redis. Also we don't need redis complex features. But we
consider port/adapter for changing the distribute cache technology.

### Local cache
Service send request to KGS for n keys and KGS response n unique generated keys. The list of keys
will divide to two parts. When half of list used, service will send independent bulk key request to KGS
and append to to list. User will not see it and it have no effect on user request response time.

*This cache strategy will reduce number of requests between shortener and KGS.*

## KGS
- Local cache for Cache pool of generated key ranges 

### Local cache
The same algorithm applied for KGS range cache but the list len will be significantly shorter for KGS.
When half of list used, service will send independent bulk key request to KGS
and append to to list. User will not see it and it have no effect on user request response time.

*This cache strategy will reduce number of requests between KGS and DB.*
