# grpc_framework

III. Using klog grpc middleware

GO Cient -- (GRPC) --> GO Server --> (GRPC) -- Python Server (Storage)

```bash
$ go run client.go getprodtypes google
I0121 09:21:43.468720   24825 client_interceptor.go:18]  "msg"="requesting all product types from vendor: google"  
I0121 09:21:43.481720   24825 client_interceptor.go:65]  "msg"="Info - The call finished with code OK"  "details"={"SystemField":"grpc client","grpc.method":"GetVendorProdTypes","grpc.service":"api.ProdService","grpc.time_ms":12.948}
google cloud products type are:  compute storage

$ go run server.go 
2020/01/21 09:21:33 Serving gRPC on https://localhost:8080
I0121 09:21:43.481232   24702 server_interceptor.go:19]  "msg"="have received a request for google as vendor "  
I0121 09:21:43.481329   24702 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-03677bff","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProdTypes","grpc.request.deadline":"2020-01-21T09:21:47+02:00","grpc.service":"api.ProdService","grpc.time_ns":15878,"peer.address":"127.0.0.1:52284"}
```

Using the storage service:

```bash
$ go run client.go getprods aws storage
I0121 09:22:46.005281   25093 client_interceptor.go:34]  "msg"="requesting all _ products from _"  
I0121 09:22:46.016896   25093 client_interceptor.go:65]  "msg"="Info - The call finished with code OK"  "details"={"SystemField":"grpc client","grpc.method":"GetVendorProds","grpc.service":"api.ProdService","grpc.time_ms":11.547}
Title: Amazon Aurora, Url: https://aws.amazon.com/rds/aurora/,  ShortUrl: https://made-up-url.com/4dcebb
Title: Amazon RDS, Url: https://aws.amazon.com/rds/,  ShortUrl: https://made-up-url.com/3a04c4
Title: Amazon Redshift, Url: https://aws.amazon.com/redshift/,  ShortUrl: https://made-up-url.com/1e989b
Title: Amazon DynamoDB, Url: https://aws.amazon.com/dynamodb/,  ShortUrl: https://made-up-url.com/e89dcb
Title: Amazon ElastiCache for Memcached, Url: https://aws.amazon.com/elasticache/memcached/,  ShortUrl: https://made-up-url.com/0e911a
Title: Amazon ElastiCache for Redis, Url: https://aws.amazon.com/elasticache/redis/,  ShortUrl: https://made-up-url.com/c697d8
Title: Amazon Neptune, Url: https://aws.amazon.com/neptune/,  ShortUrl: https://made-up-url.com/2fe1e2

$ go run server.go 
2020/01/21 09:22:39 Serving gRPC on https://localhost:8080
2020/01/21 09:22:46 have received a request for -> storage <- product type from -> aws <- vendor
I0121 09:22:46.019002   24950 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-0367aeb1","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProds","grpc.request.deadline":"2020-01-21T09:22:50+02:00","grpc.service":"api.ProdService","grpc.time_ns":1640561,"peer.address":"127.0.0.1:52290"}

$ python storage.py 
Listening on port 6000..
INFO:root:have received a request for -> storage <- product type from -> aws <- vendor
INFO:root:a number of 7 products were sent to client
```
