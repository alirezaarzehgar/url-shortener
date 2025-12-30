# Capacity Estimation

My geographical scope is limited to Iran. We hould estimate number of users, read requests,
write requests, data volume, needed RAM, needed CPU, needed servers (nodes for DB, cache, stateless apps).

Starting from number of users.
We have 90M people in Iran and 10M monthly user. Each user will create 100 short links per month.
Short links are for social media posts(Linkedin, Instagram, Twitter, etc), messages (SMS, Telegra, Whatsapp, etc)
and we approximate average 1:100 write:read ratio.

**Calculate QPS**:

```python
nusers = 10**6
monthly_wreq = 100
month = 60*60*24*30
wr_ratio = 100

write_qps = nusers * monthly_wreq / month
read_qps = write_qps * wr_ratio
```

Write ~= 40 QPS
Read ~= 4000 QPS

## Short link lenght + real link lenght

Real link lenght approximately is 1K -> **offered link lenght = 1K**
Short link lenght will be just a key.

Assuming unique keys and an arbitrary algorithm that will map data to a key.
It's depend on key encoding. We can choose encoding from base8 to base64.

We have following option:
Base10, Base16, Base20, Base32 (RFC 4648), Crockford Base32, Z-base-32
Base36, Base52, Base56, Base57, Base58, Base60, Base62
Base64url
Custom Base-N (URL-safe alphabet), Bijective Base-N (URL-safe)
Hashids (URL-safe alphabet), Sqids (URL-safe alphabet)


I choose from Base encoding. This table show number of possible link per solutions.

Calculation formula:
Base^Len.
For example we have 16^5 = 1048576 ~= 1M possible links for 5 lenght key in base16 encoding.

| Encoding          | 4 len link | 5 len   | 6 len | 7 len | 8 len |
| ----------------- | ---------- | ------- | ----- | ----- | ----- |
| Base10            | 10,000     | 100,000 | 1M    | 10M   | 100M  |
| Base16            | 65,536     | 1.05M   | 16.8M | 268M  | 4.29B |
| Base20            | 160,000    | 3.2M    | 64M   | 1.28B | 25.6B |
| Base32 (RFC 4648) | 1.05M      | 33.6M   | 1.07B | 34.4B | 1.10T |
| Crockford Base32  | 1.05M      | 33.6M   | 1.07B | 34.4B | 1.10T |
| Z-base-32         | 1.05M      | 33.6M   | 1.07B | 34.4B | 1.10T |
| Base36            | 1.68M      | 60.5M   | 2.18B | 78.4B | 2.82T |
| Base52            | 7.31M      | 380M    | 19.8B | 1.03T | 53.5T |
| Base56            | 9.83M      | 550M    | 30.8B | 1.73T | 96.6T |
| Base57            | 10.6M      | 601M    | 34.3B | 1.95T | 111T  |
| Base58            | 11.3M      | 656M    | 38.1B | 2.21T | 128T  |
| Base60            | 12.9M      | 778M    | 46.7B | 2.80T | 168T  |
| Base62            | 14.8M      | 916M    | 56.8B | 3.52T | 218T  |
| Base64url         | 16.8M      | 1.07B   | 68.7B | 4.40T | 281T  |

We can store key in various data types. Assume system will save keys
for maximum 1 year. With 40 QPS write request, system should store 1.2 billions
of unique keys.


| Storage type      | Key length | Capacity (approx) | Storage size | Total Volume (for 1.2B keys) |
|-------------------|------------|-------------------|--------------|-------------------------------|
| UINT8             | fixed      | 256               | 1 B          | X                             |
| UINT16            | fixed      | 65K               | 2 B          | X                             |
| UINT32            | fixed      | 4.3B              | 4 B          | 4.8 GB                        |
| UINT64 / BIGINT   | fixed      | 18.4E             | 8 B          | 9.6 GB                        |
| BINARY(1)         | fixed      | 256               | 1 B          | X                             |
| BINARY(2)         | fixed      | 65K               | 2 B          | X                             |
| BINARY(3)         | fixed      | 16.8M             | 3 B          | X                             |
| BINARY(4)         | fixed      | 4.3B              | 4 B          | 4.8 GB                        |
| BINARY(6)         | fixed      | 280T              | 6 B          | 7.2 GB                        |
| BINARY(8)         | fixed      | 18.4E             | 8 B          | 9.6 GB                        |
| UUID (binary)     | fixed      | 340U              | 16 B         | 19.2 GB                       |
| CHAR(4) Base16    | 4          | 65K               | 4 B          | X                             |
| CHAR(5) Base16    | 5          | 1.05M             | 5 B          | X                             |
| CHAR(6) Base16    | 6          | 16.8M             | 6 B          | X                             |
| CHAR(6) Base32    | 6          | 1.07B             | 6 B          | X                             |
| CHAR(6) Base36    | 6          | 2.18B             | 6 B          | 7.2 GB                        |
| CHAR(6) Base62    | 6          | 56.8B             | 6 B          | 7.2 GB                        |
| CHAR(7) Base62    | 7          | 3.52T             | 7 B          | 8.4 GB                        |
| CHAR(8) Base62    | 8          | 218T              | 8 B          | 9.6 GB                        |
| CHAR(5) Base64url | 5          | 1.07B             | 5 B          | X                             |
| CHAR(6) Base64url | 6          | 68.7B             | 6 B          | 7.2 GB                        |
| CHAR(7) Base64url | 7          | 4.40T             | 7 B          | 8.4 GB                        |

**Notes on large number notation:**
- **K** = Thousand (10^3)
- **M** = Million (10^6)
- **B** = Billion (10^9)
- **T** = Trillion (10^12)
- **Q** = Quadrillion (10^15)
- **E** = Quintillion (10^18)
- **D** = Decillion (10^33)
- **U** = Undecillion (10^36)

**Just keeping reasonable options**

In this scale the total volume doesn't matter and capacity is sufficient at all.
We roll out options that have no adequate capacity.

| Storage type      | Key length | Capacity (approx) | Storage size | Total Volume (for 1.2B keys) |
|-------------------|------------|-------------------|--------------|-------------------------------|
| UINT32            | fixed      | 4.3B              | 4 B          | 4.8 GB                        |
| UINT64 / BIGINT   | fixed      | 18.4E             | 8 B          | 9.6 GB                        |
| BINARY(4)         | fixed      | 4.3B              | 4 B          | 4.8 GB                        |
| BINARY(6)         | fixed      | 280T              | 6 B          | 7.2 GB                        |
| BINARY(8)         | fixed      | 18.4E             | 8 B          | 9.6 GB                        |
| UUID (binary)     | fixed      | 340U              | 16 B         | 19.2 GB                       |
| CHAR(6) Base36    | 6          | 2.18B             | 6 B          | 7.2 GB                        |
| CHAR(6) Base62    | 6          | 56.8B             | 6 B          | 7.2 GB                        |
| CHAR(7) Base62    | 7          | 3.52T             | 7 B          | 8.4 GB                        |
| CHAR(8) Base62    | 8          | 218T              | 8 B          | 9.6 GB                        |
| CHAR(6) Base64url | 6          | 68.7B             | 6 B          | 7.2 GB                        |

What is important:
- encoding
- performance
- support in databases
- complexity

I guess UINT32 is simple for developer, database, indexing and user. Also it's capacity is ~4x than what we need.

### Len

short link size = 4 bytes and encoded len = 6
offered link max len = 1K

Each link with it's information should be **~1KB**

## Storage

With 40 QPS and assuming 5 years persistent links we need **1.2 TB**

## RAM

We need cache. Probably people share short links immediately after the create them in the same day.
Each write should be cached for one day and it needs average 3.5 BG cache.

QPS * 60 * 60 * 24 / LINK SIZE

## Bandwidth

Write = 40 QPS * 1k => 40K/s
Read = 400 QPS * 1k => 4MB/s

# Conclusion

| Category             | Assumptions              | Formula / Basis       | Final Estimate              |
| -------------------- | ------------------------ | --------------------- | --------------------------- |
| **Users**            | Monthly active users     | Given                 | **10M users**               |
| **Writes**           | 100 links / user / month | Given                 | **1B links / year**         |
| **Write QPS**        | Uniform distribution     | users × links / month | **≈ 40 QPS**                |
| **Read QPS**         | 1:100 write:read ratio   | write QPS × 100       | **≈ 4,000 QPS**             |
| **Short Key**        | UINT32                   | 4 bytes               | **4 B (6 chars encoded)**   |
| **Long URL**         | Max offered length       | Given                 | **~1 KB**                   |
| **Record Size**      | URL + metadata           | Rounded               | **~1 KB / link**            |
| **Retention**        | Persistent links         | 5 years               | —                           |
| **Total Links (5y)** | 40 QPS × time            | —                     | **~1.2B links**             |
| **Storage Volume**   | 1 KB per link            | links × size          | **~1.2 TB**                 |
| **Cache Duration**   | New links                | 24 hours              | —                           |
| **Cache Size (RAM)** | One-day writes           | QPS × 86400 × 1 KB    | **~3.5 GB**                 |
| **Write Bandwidth**  | 40 QPS                   | 40 × 1 KB             | **~40 KB/s**                |
| **Read Bandwidth**   | 4,000 QPS                | 4,000 × 1 KB          | **~4 MB/s**                 |
| **Key Capacity**     | UINT32 space             | 2³²                   | **4.3B IDs (~4× headroom)** |
