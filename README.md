# grpc_framework

Example gRPC service, using Klog & Opentelemetry gRPC middleware and grpc-gateway to expose REST endpoint.

### I. GO Client

GO Cient -- (GRPC) --> GO Server --> (GRPC) -- Python Server (Storage)

#### Client -- > Server

```bash
$ go run client.go getprodtypes aws
I0127 16:32:14.155443   30549 client_interceptor.go:18]  "msg"="requesting all product types from vendor: aws"  
{
	"SpanContext": {
		"TraceID": "3e1010a186b9900d0008dc092f69444d",
		"SpanID": "d606880aedd92d3a",
		"TraceFlags": 1
	},
	"ParentSpanID": "0000000000000000",
	"SpanKind": 1,
	"Name": "grpc_tracer/Cloud-Products-types",
	"StartTime": "2020-01-27T16:32:14.155467871+02:00",
	"EndTime": "2020-01-27T16:32:14.168149486+02:00",
	"Attributes": null,
	"MessageEvents": null,
	"Links": null,
	"Status": 0,
	"HasRemoteParent": false,
	"DroppedAttributeCount": 0,
	"DroppedMessageEventCount": 0,
	"DroppedLinkCount": 0,
	"ChildSpanCount": 0
}
I0127 16:32:14.168440   30549 client_interceptor.go:65]  "msg"="Info - The call finished with code OK"  "details"={"SystemField":"grpc client","grpc.method":"GetVendorProdTypes","grpc.service":"api.ProdService","grpc.time_ms":12.941}
aws cloud products type are:  compute storage

$ go run server.go 
2020/01/27 16:31:47 Entering infinite loop
2020/01/27 16:31:47 starting HTTP/1.1 REST server on localhost:8081
2020/01/27 16:31:47 Serving gRPC on https://localhost:8080
I0127 16:32:14.167497   30315 server_interceptor.go:19]  "msg"="have received a request for aws as vendor "  
{
	"SpanContext": {
		"TraceID": "3e1010a186b9900d0008dc092f69444d",
		"SpanID": "e2e1177f993b2f04",
		"TraceFlags": 1
	},
	"ParentSpanID": "d606880aedd92d3a",
	"SpanKind": 2,
	"Name": "grpc_tracer/Cloud-Products-types",
	"StartTime": "2020-01-27T16:32:14.167556698+02:00",
	"EndTime": "2020-01-27T16:32:14.167568398+02:00",
	"Attributes": [
		{
			"Key": "grpc.server",
			"Value": {
				"Type": "STRING",
				"Value": "api-server"
			}
		}
	],
	"MessageEvents": null,
	"Links": null,
	"Status": 0,
	"HasRemoteParent": true,
	"DroppedAttributeCount": 0,
	"DroppedMessageEventCount": 0,
	"DroppedLinkCount": 0,
	"ChildSpanCount": 0
}
I0127 16:32:14.167788   30315 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-0367ee0a","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProdTypes","grpc.request.deadline":"2020-01-27T16:32:18+02:00","grpc.service":"api.ProdService","grpc.time_ns":245377,"peer.address":"127.0.0.1:55204"}
```

#### Client -- > Server -- > Storage

```bash
$ go run client.go getprods google compute
I0127 16:36:15.714252   30957 client_interceptor.go:34]  "msg"="requesting all _ products from _"  
I0127 16:36:15.725896   30957 client_interceptor.go:65]  "msg"="Info - The call finished with code OK"  "details"={"SystemField":"grpc client","grpc.method":"GetVendorProds","grpc.service":"api.ProdService","grpc.time_ms":11.575}
Title: Compute Engine, Url: https://cloud.google.com/compute/,  ShortUrl: https://made-up-url.com/2cb92b
Title: App Engine, Url: https://cloud.google.com/appengine/,  ShortUrl: https://made-up-url.com/cf559c
Title: Cloud Functions, Url: https://cloud.google.com/functions/,  ShortUrl: https://made-up-url.com/516e12
Title: Cloud Run, Url: https://cloud.google.com/run/,  ShortUrl: https://made-up-url.com/573b7c
Title: GKE, Url: https://cloud.google.com/kubernetes-engine/,  ShortUrl: https://made-up-url.com/f3acbc

$ go run server.go 
2020/01/27 16:36:52 Entering infinite loop
2020/01/27 16:36:52 starting HTTP/1.1 REST server on localhost:8081
2020/01/27 16:36:52 Serving gRPC on https://localhost:8080
2020/01/27 16:37:00 have received a request for -> compute <- product type from -> google <- vendor
I0127 16:37:00.244502   31074 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-0367b0a3","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProds","grpc.request.deadline":"2020-01-27T16:37:04+02:00","grpc.service":"api.ProdService","grpc.time_ns":1458441,"peer.address":"127.0.0.1:55276"}

$ python storage.py 
Listening on port 6000..
INFO:root:have received a request for -> compute <- product type from -> google <- vendor
INFO:root:a number of 5 products were sent to client
```

### I. REST Client

HTTP Client -- (REST) --> GO Server --> (GRPC) -- Python Server (Storage)

#### Client -- > Server

```bash
$ curl -X GET 'http://localhost:8081/api/prodtypes?vendor=oracle'
{"productType":" compute storage"}

$ go run server.go 
2020/01/27 16:41:40 Entering infinite loop
2020/01/27 16:41:40 starting HTTP/1.1 REST server on localhost:8081
2020/01/27 16:41:40 Serving gRPC on https://localhost:8080
I0127 16:42:27.447081   31394 server_interceptor.go:19]  "msg"="have received a request for oracle as vendor "  
{
	"SpanContext": {
		"TraceID": "0375fe99c6df618606054f3ae25bc5fb",
		"SpanID": "af84a7de67be7f90",
		"TraceFlags": 1
	},
	"ParentSpanID": "0000000000000000",
	"SpanKind": 2,
	"Name": "grpc_tracer/Cloud-Products-types",
	"StartTime": "2020-01-27T16:42:27.44711213+02:00",
	"EndTime": "2020-01-27T16:42:27.447127779+02:00",
	"Attributes": [
		{
			"Key": "grpc.server",
			"Value": {
				"Type": "STRING",
				"Value": "api-server"
			}
		}
	],
	"MessageEvents": null,
	"Links": null,
	"Status": 0,
	"HasRemoteParent": false,
	"DroppedAttributeCount": 0,
	"DroppedMessageEventCount": 0,
	"DroppedLinkCount": 0,
	"ChildSpanCount": 0
}
I0127 16:42:27.447348   31394 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-03674bcf","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProdTypes","grpc.service":"api.ProdService","grpc.start_time":"2020-01-27T16:42:27+02:00","grpc.time_ns":228493,"peer.address":"127.0.0.1:55346"}
```

#### Client -- > Server -- > Storage

```bash
$ curl -X GET 'http://localhost:8081/api/prods?vendor=google&productType=compute'
{"result":{"product":{"title":"Compute Engine","url":"https://cloud.google.com/compute/","shortUrl":"https://made-up-url.com/7d62d1"}}}
{"result":{"product":{"title":"App Engine","url":"https://cloud.google.com/appengine/","shortUrl":"https://made-up-url.com/ccf764"}}}
{"result":{"product":{"title":"Cloud Functions","url":"https://cloud.google.com/functions/","shortUrl":"https://made-up-url.com/302ac9"}}}
{"result":{"product":{"title":"Cloud Run","url":"https://cloud.google.com/run/","shortUrl":"https://made-up-url.com/b28804"}}}
{"result":{"product":{"title":"GKE","url":"https://cloud.google.com/kubernetes-engine/","shortUrl":"https://made-up-url.com/5ae835"}}}

$ go run server.go 
2020/01/27 16:43:57 Entering infinite loop
2020/01/27 16:43:57 starting HTTP/1.1 REST server on localhost:8081
2020/01/27 16:43:57 Serving gRPC on https://localhost:8080
2020/01/27 16:44:15 have received a request for -> compute <- product type from -> google <- vendor
I0127 16:44:15.299347   31615 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-03679b88","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProds","grpc.service":"api.ProdService","grpc.start_time":"2020-01-27T16:44:15+02:00","grpc.time_ns":1476858,"peer.address":"127.0.0.1:55376"}

$ python storage.py 
Listening on port 6000..
INFO:root:have received a request for -> compute <- product type from -> google <- vendor
INFO:root:a number of 5 products were sent to client
```