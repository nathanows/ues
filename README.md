# UES (Ultimate Echo Server)

**UES** is a sample gRPC server implementation that implements a single Echo RPC. The Echo service implements the proto contract seen below. The server itself uses TLS to secure communication and implements middleware for request lifecycle logging and Token Header authentication (preshared secret).

This repo should be seen as exemplary not productionalized. See _Implementation Details_ below for details on design decisions.

**Authentication**
Service authentication enforces passing a valid auth token on all client RPC calls, specified via an "authorization" metadata key on requests. For both client and server in this implementation this is set with a `TOKEN` env var.

---

## Test Drive

### Local

See the following functionality overview for running the application locally.

_Note: the commands below require a working Go environment and `openssl` (for generating local TLS certs)._

#### Run Server

```bash
make server       # generates local certs, compiles protos, builds binary, runs server

make serverd      # run server in Docker
```

#### Make Requests (pre-configured client)

To send a small concurrent set of preconfigured messages to the service you can use the following. For a more dynamic experience, build and use the client binary directly or use `grpcurl`

```bash
make client-request   # makes a set of concurrent sample requests to the local server
```

#### Make Requests (pre-configured client)

```bash
make build-client
TOKEN=super-secret SERVER_ADDR=localhost:6000 bin/client "message 1" "message 2"
```

#### Make Requests (gRPCurl)

`grpcurl` is a command-line tool that lets you interact with gRPC servers (think `curl` for gRPC). See [installation instructions](https://github.com/fullstorydev/grpcurl#installation) on the projects' GitHub page.

_**Note**: `--insecure` flag is needed due to usage of self-signed certs locally_

```bash
# list registered endpoints
grpcurl --insecure localhost:6000 list

# describe service
grpcurl --insecure localhost:6000 describe echo.EchoService.Echo

# describe EchoRequest
grpcurl --insecure localhost:6000 describe echo.EchoRequest

# send Echo request
grpcurl --insecure -d '{"message": "check this out"}' \
    -H 'authorization: super-secret' \
    localhost:6000 \
    echo.EchoService/Echo
```

#### Run Tests

Run the application test suite with:

```bash
make test
```

#### Compile Protos

If changes are made to the service's proto contract, proto files must be recompiled using the following command. `protoc` is run from a Docker build image so no additional installs are required.

```bash
make protogen
```

#### Build Docker Image

```bash
make container
```

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
    service.go       # implmentation of echo.EchoService RPC handlers
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

```txt
github.com/golang/protobuf
github.com/grpc-ecosystem/go-grpc-middleware
google.golang.org/grpc
```
