# p2p-grpc

fork from [cpurta/p2p-grpc](https://github.com/cpurta/p2p-grpc)

Experimental peer-to-peer network using Google's gRPC. This is a simple network whose
messages are just hello messages from the nodes to other nodes. It relies on [hashicorp
consul](https://github.com/hashicorp/consul) for service discovery of other nodes and
for node key/value storage.

## Requirements

You will need to either have consul installed locally or you can pull the docker image
so that you can have a consul agent that local nodes can connect to.

It is recommended that you also have [Go](https://golang.org/dl/) installed (1.21+)

> grpc pb generat cmd
>
> protoc -I api/pb --go_out=paths=source_relative:./api/gen/pb --go-grpc_out=paths=source_relative,require_unimplemented_servers=false:./api/gen/pb $(find api/pb/ ! -path '_google_' -name '\*.proto')

## Quick Start

You will need to build the binary by running the following command:

```
$ go build -o ./bin/p2p-grpc ./cmd/p2p-grpc
```

In another terminal tab/window you will need to start consul:

```
$ consul agent -dev
```

You can then start the first node:

```
$ ./bin/p2p-grpc --node-name node-1 --listenaddr 127.0.0.1:10000 --service-discover-addr=127.0.0.1:8500
```

and lets start another:

```
$ ./bin/p2p-grpc --node-name node-2 --listenaddr 127.0.0.1:10001 --service-discover-addr=127.0.0.1:8500
```
