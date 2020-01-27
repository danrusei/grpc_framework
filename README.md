# grpc_framework

IV. Using Opentelemetry grpc middleware

GO Cient -- (GRPC) --> GO Server --> (GRPC) -- Python Server (Storage)

```bash
$ go run client.go getprodtypes oracle
I0127 11:36:11.034955   18076 client_interceptor.go:18]  "msg"="requesting all product types from vendor: oracle"  
{
	"SpanContext": {
		"TraceID": "503ae9aae593985ae5a96823ceff7503",
		"SpanID": "515e5ceec9c60fd5",
		"TraceFlags": 1
	},
	"ParentSpanID": "0000000000000000",
	"SpanKind": 1,
	"Name": "grpc_tracer/Cloud-Products-types",
	"StartTime": "2020-01-27T11:36:11.034975934+02:00",
	"EndTime": "2020-01-27T11:36:11.047985393+02:00",
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
I0127 11:36:11.048271   18076 client_interceptor.go:65]  "msg"="Info - The call finished with code OK"  "details"={"SystemField":"grpc client","grpc.method":"GetVendorProdTypes","grpc.service":"api.ProdService","grpc.time_ms":13.261}

oracle cloud products type are:  compute storage

$ go run server.go 
2020/01/27 11:32:55 Serving gRPC on https://localhost:8080
I0127 11:36:11.047142   17698 server_interceptor.go:19]  "msg"="have received a request for oracle as vendor "  
{
	"SpanContext": {
		"TraceID": "503ae9aae593985ae5a96823ceff7503",
		"SpanID": "4780a3c993489892",
		"TraceFlags": 1
	},
	"ParentSpanID": "515e5ceec9c60fd5",
	"SpanKind": 2,
	"Name": "grpc_tracer/Cloud-Products-types",
	"StartTime": "2020-01-27T11:36:11.047211115+02:00",
	"EndTime": "2020-01-27T11:36:11.04722305+02:00",
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
I0127 11:36:11.047472   17698 server_interceptor.go:103]  "msg"="Info - finished streaming call with code OK"  "details"={"Name":"Customer-03678433","SystemField":"grpc server","grpc.code":"OK","grpc.method":"GetVendorProdTypes","grpc.request.deadline":"2020-01-27T11:36:15+02:00","grpc.service":"api.ProdService","grpc.time_ns":279834,"peer.address":"127.0.0.1:53234"}
```
