# UES (Ultimate Echo Server)

**UES** is a sample gRPC server implementation that implements a single Echo RPC. The Echo service implements the proto contract seen below. The server itself uses TLS to secure communication and implements middleware for request lifecycle logging and Token Header authentication.

This repo should be seen as exemplary not productionalized. See _Implementation Details_ below for details on design decisions.

---

## Test Drive

### Local

#### Sample Requests
To see a running demo of the functionality of this application run the following:

```
$ make server           # generates local certs, compiles protos, builds binary, runs server
$ make client-request   # makes a set of concurrent sample requests to the local server
```

The above commands require a working Go environment (for building binary), Docker (used for protobuf compilation), and `openssl` (for generating local TLS certs). To run in a fully containerized environment use the following:

```
$ make serverd           # builds and starts server container
$ make client-requestd   # builds and runs sample requests from client container
```

#### Tests

Run the application test suite with:

```
make test
```

Or, to run the tests in a containerized environment use:

```
make testd
```

### Minikube

**TODO**

### Deployed

**TODO**

---

## Implementation Details

### API Contract

The UES server exposes the EchoService which implements the following contract:

```proto
service EchoService {
  rpc Echo(EchoRequest) returns (EchoResponse);
}

message EchoRequest {
  string message = 1;
}

message EchoResponse {
  string message = 1;
}
```

### Package Layout

```sh
/bin                # compiled local binaries (not checked in)
/certs              # self-signed SSL certs for local development (not checked in)
/client             # sample API consumer
/internal           # keep it secret, keep it safe
    /pkg            # shared library abstraction for multiple services
        /middleware # gRPC interceptor middleware
/echo               # core gRPC service implmentation
    echo.pb.go      # protoc compiled, language specific proto contracts
    echo.proto      # base proto contract definition
    server.go       # implmentation of echo.EchoService RPC handlers
/server             # base gRPC server wireup
```

#### Notes

**`echo/`**

While overkill for this contrived example, extracting the service interface RPC implmentations to a separate package (gRPC in this case, but same goes for HTTP handlers) provides not only better testing surface exposure, but more importantly is structuring the server in a forward-looking, modular fashion that reduces deployment time for future services in this domain/bounded context.

If we later wanted to build a new closely related service falling within the same domain of the UES, let's say a TransformService that modified a request message (and pretend this really couldn't just be another RPC on the EchoService), we would add the `transform` package with the implementation and then wire that up in the existing server without needing to re-implement TLS, middleware, config parsing, etc..

Later on, the EchoService has turned out to be kind of a flop, but the Transform service ended up being hugely popular and the two services are scaling very differently. Due to the isolation of this package we're able to essentially fork the existing server and modify the wireup in the server and we'll have extracted that service with very little time invested.

**`internal/pkg/`**

`/internal` (details [here](https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit)) is the default directory to add most new packages to until the functionality is vetted and the API contract finalized.

`/pkg` contains isolated helper packages used across services. In this example that is just the `middleware` package, but things like DB driver wrappers, context parsers, etc. are other examples. These can generally can be good candidates for abstraction.


### Dependencies

External dependencies are managed using [Go Modules](https://github.com/golang/go/wiki/Modules), but are generally avoided where possible in this example (isolated to working with gRPC).

```
github.com/golang/protobuf
github.com/grpc-ecosystem/go-grpc-middleware
google.golang.org/grpc
```

---

## Kubernetes

**TODO**

### Cilium

**TODO**

#### EchoService Network Security Policy

**TODO**
