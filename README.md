# grpc_framework

I. Basic service

GO Cient -- (GRPC) --> GO Server --> (GRPC) -- Python Server (Storage)

```bash
$ go run client.go getprodtypes oracle
2020/01/13 11:59:44 requesting all product types from vendor: oracle
oracle cloud products type are:  compute storage

$ go run server.go 
2020/01/13 11:59:04 Serving gRPC on https://localhost:8080
2020/01/13 11:59:44 have received a request for -> oracle <- as vendor
2020/01/13 11:59:44 the response was sent to client
```

Using the storage service:

```bash
$ go run client.go getprods aws storage
2020/01/13 12:02:33 requesting all storage products from aws
Title: Amazon Aurora, Url: https://aws.amazon.com/rds/aurora/,  ShortUrl: https://made-up-url.com/2a2075
Title: Amazon RDS, Url: https://aws.amazon.com/rds/,  ShortUrl: https://made-up-url.com/6402f0
Title: Amazon Redshift, Url: https://aws.amazon.com/redshift/,  ShortUrl: https://made-up-url.com/f109d2
Title: Amazon DynamoDB, Url: https://aws.amazon.com/dynamodb/,  ShortUrl: https://made-up-url.com/6bbdbc
Title: Amazon ElastiCache for Memcached, Url: https://aws.amazon.com/elasticache/memcached/,  ShortUrl: https://made-up-url.com/1c62ae
Title: Amazon ElastiCache for Redis, Url: https://aws.amazon.com/elasticache/redis/,  ShortUrl: https://made-up-url.com/820eee
Title: Amazon Neptune, Url: https://aws.amazon.com/neptune/,  ShortUrl: https://made-up-url.com/9b42ef

$ go run server.go 
2020/01/13 12:02:22 Serving gRPC on https://localhost:8080
2020/01/13 12:02:33 have received a request for -> storage <- product type from -> aws <- vendor
2020/01/13 12:02:33 the response was sent to client

$ python storage.py 
Listening on port 6000..
INFO:root:have received a request for -> storage <- product type from -> aws <- vendor
INFO:root:a number of 7 products were sent to client
```
