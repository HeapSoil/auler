# Auler: Minimal Agent Platform Backend

A minimal golang backend for building a agent platform with basic functionalities (User Management, Prompt Management)

## Features
- Utilizes a clean and straightforward architecture.
- Employs widely-used Go packages:
    - gorm
    - govalidator
    - gin
    - cobra
    - zap
    - pprof
    - grpc
    - protobuf.
- JWT-based authentication and Casbin-based authorization.
- Independently designed log and error packages.
- HTTP, HTTPS, and gRPC servers implementations with JSON and Protobuf for data exchange.
- Supports features like call chains, graceful shutdown, middleware, CORS, and exception recovery.
-  Provides programming access to MySQL.
-  Follows RESTful API design standards.


## Installation
```
$ git clone https://github.com/HeapSoil/auler
$ go work use auler
$ cd auler
$ make
```
