CREATE KEYSPACE IF NOT EXISTS urlshortener WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'};
CREATE TABLE IF NOT EXISTS urlshortener.keypool (key int PRIMARY KEY, counter bigint);
CREATE TABLE IF NOT EXISTS urlshortener.urls (key text PRIMARY KEY, original_url text, created_at timestamp);
